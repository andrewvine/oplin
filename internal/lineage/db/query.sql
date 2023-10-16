-- name: GetJobNamespaceByID :one
select * from lineage.job_namespaces
where id = $1 limit 1;

-- name: GetJobNamespaceByName :one
select * from lineage.job_namespaces
where name = $1 limit 1;

-- name: ListJobNamespaces :many
select * from lineage.job_namespaces
order by name;

-- name: CreateJobNamespace :one
insert into lineage.job_namespaces (
  name,
  created_at
) values (
  $1, $2
)
returning *;

-- name: GetJobByID :one
select * from lineage.jobs
where id = $1 limit 1;

-- name: GetJobByNamespaceIDAndName :one
select * from lineage.jobs
where namespace_id = $1 and name = $2 limit 1;

-- name: ListJobs :many
select * from lineage.jobs
order by name;

-- name: ListJobsWithNamespaces :many
select 
  j.id, 
  j.name, 
  j.namespace_id, 
  j.facets,
  j.updated_at,
  j.created_at,
  ns.name as namespace_name,
  ns.updated_at as namespace_updated_at,
  ns.created_at as namespace_created_at
from lineage.jobs j
join lineage.job_namespaces ns on ns.id = j.namespace_id;

-- name: GetJobWithNamespace :one
select 
  j.id, 
  j.current_version_id, 
  j.name, 
  j.namespace_id, 
  j.facets,
  j.updated_at,
  j.created_at,
  ns.name as namespace_name,
  ns.updated_at as namespace_updated_at,
  ns.created_at as namespace_created_at
from lineage.jobs j
join lineage.job_namespaces ns on ns.id = j.namespace_id
where j.id = $1;

-- name: UpdateCurrentJobVersion :one
update lineage.jobs set current_version_id = $1, updated_at = $2 
where id = $3
returning *;

-- name: CreateJob :one
insert into lineage.jobs (
  namespace_id,
  name,
  facets,
  created_at
) values (
  $1, $2, $3, $4
)
returning *;


-- name: GetJobVersionByID :one
select * from lineage.job_versions
where id = $1 limit 1;

-- name: ListJobVersionsByJobID :many
select * from lineage.job_versions
where job_id = $1 
order by created_at desc;

-- name: CreateJobVersion :one
insert into lineage.job_versions (
  job_id,
  namespace_id,
  name,
  facets,
  created_at
) values (
  $1, $2, $3, $4, $5
)
returning *;

-- name: GetRunByID :one
select * from lineage.runs
where id = $1 limit 1;


-- name: GetRunByUUID :one
select * from lineage.runs
where run_uuid = $1 limit 1;

-- name: ListRuns :many
select * from lineage.runs
order by job_version_id, id;

-- name: ListRunsByJobVersionID :many
select * from lineage.runs
where job_version_id = $1
order by created_at desc;

-- name: CreateRun :one
INSERT INTO lineage.runs (
  run_uuid,
  job_version_id,
  facets,
  parent_run_id,
  started_at,
  ended_at,
  nominal_started_at,
  nominal_ended_at,
  error_message,
  programming_language,
  stacktrace,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateRun :one
UPDATE lineage.runs SET 
  facets = $2,
  last_event_type = $3,
  error_message = $4,
  programming_language = $5,
  stacktrace = $6,
  ended_at = $7,
  updated_at = $8
WHERE id = $1
RETURNING *;

-- name: GetRunEvent :one
SELECT * FROM lineage.run_events
WHERE id = $1 LIMIT 1;

-- name: ListRunEventsByRunID :many
SELECT * FROM lineage.run_events
WHERE run_id = $1
ORDER BY created_at asc;

-- name: CreateRunEvent :one
INSERT INTO lineage.run_events (
  run_id,
  event_type,
  event_time,
  facets,
  created_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetDatasetNamespaceByID :one
select * from lineage.dataset_namespaces
where id = $1 limit 1;

-- name: GetDatasetNamespaceByName :one
select * from lineage.dataset_namespaces
where name = $1 limit 1;

-- name: ListDatasetNamespaces :many
select * from lineage.dataset_namespaces
order by name;

-- name: CreateDatasetNamespace :one
insert into lineage.dataset_namespaces (
  name,
  created_at
) values (
  $1, $2
)
returning *;


-- name: GetDatasetByID :one
select * from lineage.datasets
where id = $1 limit 1;

-- name: GetDatasetByNamespaceIDAndName :one
select * from lineage.datasets
where namespace_id = $1 and name = $2 limit 1;

-- name: CreateDataset :one
insert into lineage.datasets (
  namespace_id,
  name,
  facets,
  created_at
) values (
  $1, $2, $3, $4
)
returning *;


-- name: ListDatasetsWithNamespaces :many
select 
  d.id, 
  d.current_version_id,
  d.name, 
  d.namespace_id, 
  d.facets,
  d.updated_at,
  d.created_at,
  ns.name as namespace_name,
  ns.updated_at as namespace_updated_at,
  ns.created_at as namespace_created_at
from lineage.datasets d
join lineage.dataset_namespaces ns on ns.id = d.namespace_id;

-- name: GetDatasetWithNamespace :one
select 
  d.id, 
  d.current_version_id,
  d.name, 
  d.namespace_id, 
  d.facets,
  d.updated_at,
  d.created_at,
  ns.name as namespace_name,
  ns.updated_at as namespace_updated_at,
  ns.created_at as namespace_created_at
from lineage.datasets d
join lineage.dataset_namespaces ns on ns.id = d.namespace_id
where d.id = $1;


-- name: UpdateCurrentDatasetVersion :one
update lineage.datasets set current_version_id = $1, updated_at = $2 
where id = $3
returning *;

-- name: UpdateDataset :one
update lineage.datasets set 
  facets = $1,
  updated_at = $2 
where id = $3
returning *;

-- name: GetDatasetVersionByID :one
select * from lineage.dataset_versions
where id = $1 limit 1;

-- name: ListDatasetVersionsByDatasetID :many
select * from lineage.dataset_versions
where dataset_id = $1 
order by created_at desc;

-- name: CreateDatasetVersion :one
insert into lineage.dataset_versions (
  dataset_id,
  namespace_id,
  name,
  created_at
) values (
  $1, $2, $3, $4
)
returning *;

-- name: GetRunDatasetVersionByRunIDAndDatasetVersionID :one
select * from lineage.run_dataset_versions
where run_id = $1 and dataset_version_id = $2 limit 1;

-- name: CreateRunDatasetVersion :one
insert into lineage.run_dataset_versions (
  run_id,
  dataset_version_id,
  io_type,
  dataset_facets,
  io_facets,
  created_at
) values (
  $1, $2, $3, $4, $5, $6
)
returning *;

-- name: GetLatestRunDatasetVersionByDatasetVersionID :one
select * from lineage.run_dataset_versions
where dataset_version_id = $1 order by created_at desc limit 1;

-- name: ListRunDatasetVersionsWithRelationshipsByRunID :many
select 
  r.run_id, 
  r.dataset_version_id,
  r.io_type,
  r.dataset_facets,
  r.io_facets,
  r.created_at,
  v.id as version_id, 
  v.dataset_id as version_dataset_id, 
  v.name as version_name, 
  v.namespace_id as version_namespace_id, 
  v.updated_at as version_updated_at,
  v.created_at as version_created_at,
  n.id as namespace_id, 
  n.name as namespace_name,
  n.updated_at as namespace_updated_at,
  n.created_at as namespace_created_at
from lineage.run_dataset_versions r
join lineage.dataset_versions v on v.id = r.dataset_version_id
join lineage.dataset_namespaces n on n.id = v.namespace_id
where r.run_id = $1;

-- name: CreateField :one
insert into lineage.fields (
  dataset_version_id,
  name,
  data_type,
  description,
  updated_at,
  created_at
) values (
  $1, $2, $3, $4, $5, $6
)
returning *;

-- name: ListFieldsByDatasetVersionID :many
select * from lineage.fields
where dataset_version_id = $1 order by name; 

-- name: CreateLifecycleStateChange :one
insert into lineage.lifecycle_state_changes (
  dataset_id,
  namespace,
  name,
  data_type,
  change,
  updated_at,
  created_at
) values (
  $1, $2, $3, $4, $5, $6, $7
)
returning *;

-- name: ListLifecycleStateChangesByDatasetID :many
select * from lineage.lifecycle_state_changes
where dataset_id = $1 order by created_at; 

-- name: CreateRequest :one
INSERT INTO lineage.requests (
  payload,
  created_at
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListRequests :many
select * from lineage.requests
order by created_at; 