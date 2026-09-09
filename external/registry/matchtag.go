package registry

import (
	"fmt"
	"mime"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var versionTagPattern = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)

type matchManifest struct {
	SchemaVersion int          `json:"schemaVersion"`
	MediaType     string       `json:"mediaType"`
	Config        ManifestInfo `json:"config"`
	FSLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
}

type versionTag struct {
	tag                 string
	major, minor, patch uint64
	prerelease          []string
}

type manifestIdentity struct {
	configDigest string
	layerDigests string
}

// MatchTag returns the highest semantic-version tag whose image content matches reference.
func (r *RegistryClient) MatchTag(repository, reference string) (string, error) {
	tags, _, err := r.GetTags(repository)
	if err != nil {
		return "", fmt.Errorf("error getting tags for %s: %v", repository, err)
	}

	referenceIdentity, err := r.manifestIdentity(repository, reference)
	if err != nil {
		return "", fmt.Errorf("error resolving %s %s: %v", repository, reference, err)
	}

	var highest versionTag
	found := false
	for _, tag := range tags.Tags {
		version, ok := parseVersionTag(tag)
		if !ok {
			continue
		}
		identity, err := r.manifestIdentity(repository, tag)
		if err != nil {
			return "", fmt.Errorf("error resolving %s %s: %v", repository, tag, err)
		}
		if identity == referenceIdentity && (!found || version.greaterThan(highest)) {
			highest = version
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("no versioned tag matches %s %s", repository, reference)
	}
	return highest.tag, nil
}

func (r *RegistryClient) manifestIdentity(repository, reference string) (manifestIdentity, error) {
	u := fmt.Sprintf(r.baseUrl+manifestsPath, repository, reference)
	h := r.GetCustomHeader(fmt.Sprintf("%s, %s, %s, %s", MIME_V2_MANIFEST, MIME_OCI_MANIFEST, MIME_V1_MANIFEST, MIME_V1_PRETTYJWS))
	manifest, responseHeaders, err := HttpDo[matchManifest](r.httpClient, http.MethodGet, u, h, nil)
	if err != nil {
		return manifestIdentity{}, err
	}
	mediaType := manifest.MediaType
	if mediaType == "" {
		mediaType, _, _ = mime.ParseMediaType(responseHeaders.Get("Content-Type"))
	}
	switch mediaType {
	case MIME_V2_MANIFEST, MIME_OCI_MANIFEST:
		if manifest.Config.Digest == "" {
			return manifestIdentity{}, fmt.Errorf("manifest %s has no config digest", reference)
		}
		return manifestIdentity{configDigest: manifest.Config.Digest}, nil
	case MIME_V1_MANIFEST, MIME_V1_PRETTYJWS:
		if len(manifest.FSLayers) == 0 {
			return manifestIdentity{}, fmt.Errorf("schema 1 manifest %s has no layers", reference)
		}
		layers := make([]string, len(manifest.FSLayers))
		for i, layer := range manifest.FSLayers {
			layers[i] = layer.BlobSum
		}
		sort.Strings(layers)
		return manifestIdentity{layerDigests: strings.Join(layers, "\x00")}, nil
	default:
		return manifestIdentity{}, fmt.Errorf("unsupported manifest media type %q", mediaType)
	}
}

func parseVersionTag(tag string) (versionTag, bool) {
	parts := versionTagPattern.FindStringSubmatch(tag)
	if parts == nil {
		return versionTag{}, false
	}
	values := make([]uint64, 3)
	for i := range values {
		value, err := strconv.ParseUint(parts[i+1], 10, 64)
		if err != nil {
			return versionTag{}, false
		}
		values[i] = value
	}
	var prerelease []string
	if parts[4] != "" {
		prerelease = strings.Split(parts[4], ".")
		for _, identifier := range prerelease {
			if isNumericIdentifier(identifier) && len(identifier) > 1 && identifier[0] == '0' {
				return versionTag{}, false
			}
		}
	}
	return versionTag{tag: tag, major: values[0], minor: values[1], patch: values[2], prerelease: prerelease}, true
}

func (v versionTag) greaterThan(other versionTag) bool {
	if v.major != other.major {
		return v.major > other.major
	}
	if v.minor != other.minor {
		return v.minor > other.minor
	}
	if v.patch != other.patch {
		return v.patch > other.patch
	}
	return comparePrerelease(v.prerelease, other.prerelease) > 0
}

func comparePrerelease(left, right []string) int {
	if len(left) == 0 {
		if len(right) == 0 {
			return 0
		}
		return 1
	}
	if len(right) == 0 {
		return -1
	}
	for i := 0; i < len(left) && i < len(right); i++ {
		if left[i] == right[i] {
			continue
		}
		leftNumeric := isNumericIdentifier(left[i])
		rightNumeric := isNumericIdentifier(right[i])
		switch {
		case leftNumeric && rightNumeric:
			if len(left[i]) != len(right[i]) {
				if len(left[i]) > len(right[i]) {
					return 1
				}
				return -1
			}
		case leftNumeric:
			return -1
		case rightNumeric:
			return 1
		}
		if left[i] > right[i] {
			return 1
		}
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	if len(left) < len(right) {
		return -1
	}
	return 0
}

func isNumericIdentifier(identifier string) bool {
	for _, character := range identifier {
		if character < '0' || character > '9' {
			return false
		}
	}
	return identifier != ""
}
