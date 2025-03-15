package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// type BlobsResponse struct {
//   "architecture": "amd64",
//   "config": {
//     "User": "tomcat",
//     "ExposedPorts": {
//       "8009/tcp": {},
//       "8080/tcp": {}
//     },
//     "Env": [
//       "PATH=/opt/bitnami/common/bin:/opt/bitnami/java/bin:/opt/bitnami/tomcat/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
//       "HOME=/",
//       "OS_ARCH=amd64",
//       "OS_FLAVOUR=debian-12",
//       "OS_NAME=linux",
//       "APP_VERSION=10.1.34",
//       "BITNAMI_APP_NAME=tomcat",
//       "JAVA_HOME=/opt/bitnami/java",
//       "JASPER_VERSION=8.2.0",
//       "JASPER_HOME=/opt/bitnami/tomcat/webapps/jasperserver",
//       "PGSQL_JDBC_VERSION=42.5.0",
//       "CLICKHOUSE_JDBC_VERSION=0.8.0"
//     ],
//     "Entrypoint": [
//       "/opt/bitnami/scripts/tomcat/entrypoint.sh"
//     ],
//     "Cmd": [
//       "/opt/bitnami/scripts/tomcat/run.sh"
//     ],
//     "WorkingDir": "/opt/bitnami/tomcat/webapps/jasperserver",
//     "Labels": {
//       "com.vmware.cp.artifact.flavor": "sha256:c50c90cfd9d12b445b011e6ad529f1ad3daea45c26d20b00732fae3cd71f6a83",
//       "org.opencontainers.image.base.name": "docker.io/bitnami/minideb:bookworm",
//       "org.opencontainers.image.created": "2025-02-11T20:49:58Z",
//       "org.opencontainers.image.description": "Application packaged by Broadcom, Inc.",
//       "org.opencontainers.image.documentation": "https://github.com/bitnami/containers/tree/main/bitnami/tomcat/README.md",
//       "org.opencontainers.image.licenses": "Apache-2.0",
//       "org.opencontainers.image.ref.name": "10.1.34-debian-12-r5",
//       "org.opencontainers.image.source": "https://github.com/bitnami/containers/tree/main/bitnami/tomcat",
//       "org.opencontainers.image.title": "tomcat",
//       "org.opencontainers.image.vendor": "Broadcom, Inc.",
//       "org.opencontainers.image.version": "10.1.34"
//     },
//     "ArgsEscaped": true,
//     "Shell": [
//       "/bin/bash",
//       "-o",
//       "errexit",
//       "-o",
//       "nounset",
//       "-o",
//       "pipefail",
//       "-c"
//     ]
//   },
//   "created": "2025-02-13T14:01:52.814525715+11:00",
//   "history": [
//     {
//       "created": "0001-01-01T00:00:00Z",
//       "created_by": "crane flatten sha256:d10e085f6c6b162693465c347135f038ea0e23afe81ca47f2138cf4312768485",
//       "comment": "[{\"created\":\"2025-02-11T07:43:24.594606925Z\",\"comment\":\"from Bitnami with love\"},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"ARG DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"ARG JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"ARG TARGETARCH=amd64\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"LABEL com.vmware.cp.artifact.flavor=sha256:c50c90cfd9d12b445b011e6ad529f1ad3daea45c26d20b00732fae3cd71f6a83 org.opencontainers.image.base.name=docker.io/bitnami/minideb:bookworm org.opencontainers.image.created=2025-02-11T20:49:58Z org.opencontainers.image.description=Application packaged by Broadcom, Inc. org.opencontainers.image.documentation=https://github.com/bitnami/containers/tree/main/bitnami/tomcat/README.md org.opencontainers.image.licenses=Apache-2.0 org.opencontainers.image.ref.name=10.1.34-debian-12-r5 org.opencontainers.image.source=https://github.com/bitnami/containers/tree/main/bitnami/tomcat org.opencontainers.image.title=tomcat org.opencontainers.image.vendor=Broadcom, Inc. org.opencontainers.image.version=10.1.34\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"ENV HOME=/ OS_ARCH=amd64 OS_FLAVOUR=debian-12 OS_NAME=linux\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"COPY prebuildfs / # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:03.774185116Z\",\"created_by\":\"SHELL [/bin/bash -o errexit -o nounset -o pipefail -c]\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:10.409829411Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c install_packages ca-certificates curl libssl3 procps xmlstarlet zlib1g # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:12.80082222Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c mkdir -p /tmp/bitnami/pkg/cache/ ; cd /tmp/bitnami/pkg/cache/ ;     COMPONENTS=(       \\\"render-template-1.0.7-11-linux-${OS_ARCH}-debian-12\\\"       \\\"jre-21.0.6-10-0-linux-${OS_ARCH}-debian-12\\\"       \\\"tomcat-10.1.34-2-linux-${OS_ARCH}-debian-12\\\"     ) ;     for COMPONENT in \\\"${COMPONENTS[@]}\\\"; do       if [ ! -f \\\"${COMPONENT}.tar.gz\\\" ]; then         curl -SsLf \\\"https://${DOWNLOADS_URL}/${COMPONENT}.tar.gz\\\" -O ;         curl -SsLf \\\"https://${DOWNLOADS_URL}/${COMPONENT}.tar.gz.sha256\\\" -O ;       fi ;       sha256sum -c \\\"${COMPONENT}.tar.gz.sha256\\\" ;       tar -zxf \\\"${COMPONENT}.tar.gz\\\" -C /opt/bitnami --strip-components=2 --no-same-owner --wildcards '*/files' ;       rm -rf \\\"${COMPONENT}\\\".tar.gz{,.sha256} ;     done # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:15.507912352Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c apt-get autoremove --purge -y curl \\u0026\\u0026     apt-get update \\u0026\\u0026 apt-get upgrade -y \\u0026\\u0026     apt-get clean \\u0026\\u0026 rm -rf /var/lib/apt/lists /var/cache/apt/archives # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:15.996775433Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c chmod g+rwX /opt/bitnami # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:16.571816913Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c find / -perm /6000 -type f -exec chmod a-s {} \\\\; || true # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:17.041343714Z\",\"created_by\":\"COPY rootfs / # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:17.543082387Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c /opt/bitnami/scripts/java/postunpack.sh # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"RUN |3 DOWNLOADS_URL=dye1tjwtyxcc2.cloudfront.net/tac-rel/nami-components JAVA_EXTRA_SECURITY_DIR=/bitnami/java/extra-security TARGETARCH=amd64 /bin/bash -o errexit -o nounset -o pipefail -c /opt/bitnami/scripts/tomcat/postunpack.sh # buildkit\",\"comment\":\"buildkit.dockerfile.v0\"},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"ENV APP_VERSION=10.1.34 BITNAMI_APP_NAME=tomcat JAVA_HOME=/opt/bitnami/java PATH=/opt/bitnami/common/bin:/opt/bitnami/java/bin:/opt/bitnami/tomcat/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"EXPOSE map[8009/tcp:{} 8080/tcp:{}]\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"USER 1001\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"ENTRYPOINT [\\\"/opt/bitnami/scripts/tomcat/entrypoint.sh\\\"]\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true},{\"created\":\"2025-02-11T20:51:18.917023749Z\",\"created_by\":\"CMD [\\\"/opt/bitnami/scripts/tomcat/run.sh\\\"]\",\"comment\":\"buildkit.dockerfile.v0\",\"empty_layer\":true}]"
//     },
//     {
//       "created": "2025-02-13T11:54:50.161779526+11:00",
//       "created_by": "USER root",
//       "comment": "buildkit.dockerfile.v0",
//       "empty_layer": true
//     },
//     {
//       "created": "2025-02-13T11:54:50.161779526+11:00",
//       "created_by": "RUN /bin/bash -o errexit -o nounset -o pipefail -c apt-get update && apt-get install -y unzip wget && rm -rf /var/lib/apt/lists/* # buildkit",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T11:54:50.161779526+11:00",
//       "created_by": "ENV JASPER_VERSION=8.2.0 JASPER_HOME=/opt/bitnami/tomcat/webapps/jasperserver PGSQL_JDBC_VERSION=42.5.0 CLICKHOUSE_JDBC_VERSION=0.8.0",
//       "comment": "buildkit.dockerfile.v0",
//       "empty_layer": true
//     },
//     {
//       "created": "2025-02-13T11:54:50.517028069+11:00",
//       "created_by": "WORKDIR /opt/bitnami/tomcat/webapps/jasperserver",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T12:00:57.723281196+11:00",
//       "created_by": "COPY /tmp/jasperreports-server-cp-8.2.0-bin/* . # buildkit",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T12:02:02.93285261+11:00",
//       "created_by": "RUN /bin/bash -o errexit -o nounset -o pipefail -c chown -R tomcat:tomcat ${JASPER_HOME} # buildkit",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T14:01:20.382573897+11:00",
//       "created_by": "RUN /bin/bash -o errexit -o nounset -o pipefail -c wget -q https://github.com/ClickHouse/clickhouse-java/releases/download/v${CLICKHOUSE_JDBC_VERSION}/clickhouse-jdbc-${CLICKHOUSE_JDBC_VERSION}-all.jar   -O /opt/bitnami/tomcat/lib/clickhouse-jdbc-${CLICKHOUSE_JDBC_VERSION}-all.jar   && chmod 644 /opt/bitnami/tomcat/lib/clickhouse-jdbc-${CLICKHOUSE_JDBC_VERSION}-all.jar # buildkit",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T14:01:52.814525715+11:00",
//       "created_by": "RUN /bin/bash -o errexit -o nounset -o pipefail -c chown -R tomcat:tomcat ${JASPER_HOME} # buildkit",
//       "comment": "buildkit.dockerfile.v0"
//     },
//     {
//       "created": "2025-02-13T14:01:52.814525715+11:00",
//       "created_by": "EXPOSE map[8080/tcp:{}]",
//       "comment": "buildkit.dockerfile.v0",
//       "empty_layer": true
//     },
//     {
//       "created": "2025-02-13T14:01:52.814525715+11:00",
//       "created_by": "USER tomcat",
//       "comment": "buildkit.dockerfile.v0",
//       "empty_layer": true
//     },
//     {
//       "created": "2025-02-13T14:01:52.814525715+11:00",
//       "created_by": "CMD [\"/opt/bitnami/scripts/tomcat/run.sh\"]",
//       "comment": "buildkit.dockerfile.v0",
//       "empty_layer": true
//     }
//   ],
//   "os": "linux",
//   "rootfs": {
//     "type": "layers",
//     "diff_ids": [
//       "sha256:c173ffe16c66572a7bddeffa862cd41fb17267b4cca13c52b37256a98b485745",
//       "sha256:c2bdc9bbad402a48fbd6551d6bb8a0ba08ba49f71a81818fd0a289bc661b2046",
//       "sha256:70db68a6c2a09e7ab0b0890a21fdc0186fa6805c2ea6795becde44a57b928ed6",
//       "sha256:4360a25537a3c470eaf3d9e21885cd4236022178c373fd9942235ebf63502d00",
//       "sha256:3938fcb4e9a9217a55347f97bd6e27cd3e3227336eeb803914b1cce10c83f3c9",
//       "sha256:8d604134fc8b8765239a6aae3d6ffc2ec03fbfdc9263b6a7c157676427640c89",
//       "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
//     ]
//   }
// }

const (
	tagsPath      = `%s/tags/list`
	manifestsPath = `%s/manifests/%s`
	blobsPath     = `%s/blobs/%s`
)

// Response Get Tags
// -----------------
//
// 200 OK
// Content-Type: application/json
//
//	{
//	    "name": <name>,
//	    "tags": [
//	        <tag>,
//	        ...
//	    ]
//	}
func (r *Registry) GetTags(repository string) ([]string, error) {
	var tags []string
	u := fmt.Sprintf(r.baseUrl+tagsPath, repository)
	response, _, err := HttpDo[TagsResponse](r.httpClient, http.MethodGet, u, r.httpHeaders, nil)
	if err != nil {
		return tags, fmt.Errorf("error gettings tags: %v", err)
	}
	tags = response.Tags
	return tags, nil
}

// WARNING: quetion remains about the 2 firsts 1st mediaType to use,
// list tags : index, and manifest v1 works OR index, and manifests list v2
// manifests tag : index v1
// manifests "platform" digest : use provided mediaType info
// blobs config : use provided mediaType info
// 💡 include all headers when unknown : index v1 + manifests  list v2
// TODO : index  manifest, then image manifest with blobs, then ?!.... to retrieve CreatedAt info
// 1) get manifests/<reference>  with application/vnd.oci.image.index.v1+json
// HTTP/1.1 200 OK
// Content-Length: 856
// Content-Type: application/vnd.oci.image.index.v1+json
// Docker-Content-Digest: sha256:66339d5454990689da17ffb99217e333c9d84600344d09ee758cff8f6594cb90
// Docker-Distribution-Api-Version: registry/2.0
// Etag: "sha256:66339d5454990689da17ffb99217e333c9d84600344d09ee758cff8f6594cb90"
// X-Content-Type-Options: nosniff
// Date: Mon, 17 Feb 2025 20:46:53 GMT
//
//	{
//	  "schemaVersion": 2,
//	  "mediaType": "application/vnd.oci.image.index.v1+json",
//	  "manifests": [
//	    {
//	      "mediaType": "application/vnd.oci.image.manifest.v1+json",
//	      "digest": "sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17",
//	      "size": 1631,
//	      "platform": {
//	        "architecture": "amd64",
//	        "os": "linux"
//	      }
//	    },
//	    {
//	      "mediaType": "application/vnd.oci.image.manifest.v1+json",
//	      "digest": "sha256:02d363a44f31dd0def7dcccff9fbe3b09b2e806cb44aee3e03baf104d097c5a6",
//	      "size": 566,
//	      "annotations": {
//	        "vnd.docker.reference.digest": "sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17",
//	        "vnd.docker.reference.type": "attestation-manifest"
//	      },
//	      "platform": {
//	        "architecture": "unknown",
//	        "os": "unknown"
//	      }
//	    }
//	  ]
//	}%
//
// 2) get manifests sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17 with application/vnd.oci.image.manifest.v1+json
// {
// "schemaVersion": 2,
// "mediaType": "application/vnd.oci.image.manifest.v1+json",
//
//	"config": {
//	  "mediaType": "application/vnd.oci.image.config.v1+json",
//	  "digest": "sha256:1649f157365545ac4b8ec167619fb18d2b61f802776e39e46a8156f39762615e",
//	  "size": 11148
//	},
//
// "layers": [
//
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278",
//	  "size": 145483495
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:b9be66bfe7f92b5c42a47c6353d0dfb1f7b9610a9479752228d5f1fe00c100fc",
//	  "size": 2094433
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:08d8d343d6a4c6fb7033d42667d66a88368e95b0f1ee288621dbaf24149d33ca",
//	  "size": 178
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:e6c0e3d5828e19ef46e585f10e2af75e11be87f42432301fb92df598d2d2d092",
//	  "size": 477195673
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:e8ff69f6858575d6e0a8be832b30716e41bce7379acee914fa91da26533a484a",
//	  "size": 477197575
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:107aba61455803961e2bf3981ee15312ff09177dfa623fe785f4759b63afa9a5",
//	  "size": 7267273
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1",
//	  "size": 32
//	}
//
// ]
//
// 3) get blobs/sha256:1649f157365545ac4b8ec167619fb18d2b61f802776e39e46a8156f39762615e with application/vnd.oci.image.config.v1+json
// cf blobs.sample.json
func (r *Registry) Inspect(name, reference string) (string, error) {
	u := fmt.Sprintf(r.baseUrl+manifestPath, name, reference)

	req, err := GetNewRequest(http.MethodGet, u, r.httpHeaders, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting manifest for %s %s: %v", name, reference, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var respErr RegistryError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return "", fmt.Errorf("%d, error getting manifest: %v", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("%d, error getting manifest: %v", resp.StatusCode, respErr)
	}

	return "", nil
}
