# helmsplain

Show templated values in files

```bash
$ helm pull bitnami/postgresql --version 12.1.1
$ ./bin/helmsplain ~/Downloads/postgresql-12.1.1.tgz
/postgresql/charts/common/README.md
     .Values.password
/postgresql/templates/primary/statefulset.yaml
     .Values.primary.hostNetwork
     .Values.primary.hostIPC
     .Values.image.pullPolicy
     .Values.containerPorts.postgresql
     .Values.primary.persistence.mountPath
     .Values.audit.logHostname
     .Values.audit.logConnections
     .Values.audit.logDisconnections
     .Values.audit.pgAuditLogCatalog
     .Values.audit.clientMinMessages
     .Values.postgresqlSharedPreloadLibraries
     .Values.containerPorts.postgresql
/postgresql/templates/primary/svc.yaml
     .Values.primary.service.type
/postgresql/values.yaml
     .Values.metrics.service.ports.metrics
```