Minion
=======

This system runs on all edge & origins machines and provides the proxy files to the servers.
It also handles a bunch of utility features such as file uploading and validation.

# Artifact

Access to the compiled binary is available from the following link. Remeber to replace `buildid`
http://deployer:eZKpb4uR9cEW9A7tvQF2ug3LCNZxxpzysuDYYBzbUn@teamcity.sla.bz/repository/download/Minion_Build/:buildid/minion.zip


## Http endpoints

The minion implements two endpoints, one for downloading and proxy and one for uploading new content

### Download

Endpoint: `/download`

Method: *get*

Query parameters

| Name | Purpose |
| ---- | ------- |
| host | Target host to download the file from |
| path | Target path of the file from which we want to download  |

### Upload

Endpoint `/upload`
Method: *post*

Form parameters

| Name | purpose |
| ---- | ------- |
| File | the actual file being uploaded |
| ExternalId | the external id of the file being uploaded uint64 |
| FileType | Check the protobufer bindings for this |
