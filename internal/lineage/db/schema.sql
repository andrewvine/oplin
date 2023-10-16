create schema lineage;

create table lineage.job_namespaces (
  id             bigserial primary key,
  name           varchar(255) not null,
  created_at     timestamp not null, 
  updated_at     timestamp,
  unique(name)
);

create table lineage.jobs (
  id                       bigserial primary key,
  current_version_id       bigint,
  namespace_id   bigint not null,
  name           varchar(255) not null,
  facets         jsonb,
  created_at     timestamp not null, 
  updated_at     timestamp,
  unique(namespace_id, name),
  constraint     
    fk_namespace foreign key(namespace_id) 
      references lineage.job_namespaces(id)
);

create table lineage.job_versions (
  id             bigserial primary key,
  job_id         bigint  not null,
  namespace_id   bigint  not null,
  name           varchar(255) not null,
  facets         jsonb,
  created_at     timestamp not null, 
  updated_at     timestamp, 
  constraint     
    fk_job_id foreign key(job_id) 
      references lineage.jobs(id),
  constraint     
    fk_namespace foreign key(namespace_id) 
      references lineage.job_namespaces(id)
);

create table lineage.runs (
  id                    bigserial primary key,
  run_uuid              uuid      not null,
  job_version_id        bigint    not null,
  parent_run_id         bigint,
  last_event_type       int not null default 0,
  facets                jsonb,
  started_at            timestamp,
  ended_at              timestamp, 
  nominal_started_at    timestamp,
  nominal_ended_at      timestamp, 
  error_message         varchar, 
  programming_language  varchar, 
  stacktrace            varchar, 
  created_at            timestamp not null,
  updated_at            timestamp, 
  unique(run_uuid),
  constraint     
    fk_job_version_id foreign key(job_version_id) 
      references lineage.job_versions(id)
);

create table lineage.run_events (
  id              bigserial    primary key,
  run_id          bigint       not null,
  event_type      int          not null default 0,
  event_time      timestamp    not null,
  facets          jsonb,
  created_at      timestamp not null, 
  updated_at      timestamp, 
  constraint     
    fk_run foreign key(run_id) 
      references lineage.runs(id)
);

create table lineage.dataset_namespaces (
  id             bigserial primary key,
  name           varchar(255) not null,
  created_at     timestamp not null, 
  updated_at     timestamp, 
  unique(name)
);

create table lineage.datasets (
  id                        bigserial primary key,
  current_version_id        bigint,
  namespace_id              bigint not null,
  name                      varchar(255) not null,
  facets                    jsonb,
  created_at                timestamp not null, 
  updated_at                timestamp,
  unique(namespace_id, name),
  constraint     
    fk_namespace foreign key(namespace_id) 
      references lineage.dataset_namespaces(id)
);

create table lineage.dataset_versions (
  id             bigserial primary key,
  dataset_id     bigint  not null,
  namespace_id   bigint not null,
  name           varchar(255) not null,
  created_at     timestamp not null, 
  updated_at     timestamp,
  constraint     
    fk_dataset_id foreign key(dataset_id) 
      references lineage.datasets(id),
  constraint     
    fk_namespace foreign key(namespace_id) 
      references lineage.dataset_namespaces(id)
);

create table lineage.run_dataset_versions (
  run_id                 bigint not null,
  dataset_version_id     bigint not null,
  io_type                int not null, -- INPUT|OUTPUT
  dataset_facets         jsonb,
  io_facets              jsonb,
  created_at             timestamp not null,
  primary key (run_id, dataset_version_id),
  constraint     
    fk_dataset_version_id foreign key(dataset_version_id) 
      references lineage.dataset_versions(id),
  constraint     
    fk_run_id foreign key(run_id) 
      references lineage.runs(id)
);

create table lineage.fields (
  id                            bigserial primary key,
  dataset_version_id            bigint not null,
  name                          varchar not null,
  data_type                     varchar not null,
  description                   varchar,
  created_at                    timestamp not null,
  updated_at                    timestamp,
  unique(dataset_version_id, name),
  constraint     
    fk_dataset_version_id foreign key(dataset_version_id) 
      references lineage.dataset_versions(id)
);

create table lineage.lifecycle_state_changes (
  id                            bigserial primary key,
  dataset_id                    bigint not null,
  change                        varchar,
  namespace                     varchar,
  name                          varchar,
  data_type                     varchar not null,
  created_at                    timestamp not null,
  updated_at                    timestamp,
  constraint     
    fk_dataset_id foreign key(dataset_id) 
      references lineage.datasets(id)
);

create table lineage.requests (
  id              bigserial primary key,
  payload         jsonb not null,
  created_at      timestamp not null
);
