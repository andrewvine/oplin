from uuid import uuid4

from openlineage.client import client, facet, run, set_producer
from openlineage.client.run import RunState
from openlineage.client.serde import Serde


LineageFields = facet.ColumnLineageDatasetFacetFieldsAdditional
InputFields = facet.ColumnLineageDatasetFacetFieldsAdditionalInputFields
Identifiers = facet.SymlinksDatasetFacetIdentifiers
Owners = facet.OwnershipDatasetFacetOwners

source_brands = run.InputDataset(
    namespace="retail_source",
    name="brands",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Brand identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of brand"
                ),
                facet.SchemaField(
                    name="country", type="varchar", description="country code"
                )
            ]
        ),
        "dataSource": facet.DataSourceDatasetFacet(
            name="Retail RDS database", uri="postgresql://dataops@rds/retail"
        ),
    },
    inputFacets=facet.DataQualityMetricsInputDatasetFacet(
        rowCount=1310, bytes=28929201, columnMetrics={
            "id": facet.ColumnMetric(distinctCount=1210)
        }
    )
)

source_products = run.Dataset(
    namespace="retail_source",
    name="products",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Product identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of product"
                ),
                facet.SchemaField(
                    name="price", type="bigdecimal", description="Price of product"
                ),
                facet.SchemaField(
                    name="brand_id", type="varchar", description="Brand identifier"
                ),
            ]
        ),
    },
)


source_sales = run.Dataset(
    namespace="retail_source",
    name="sales",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Sale identifier"
                ),
                facet.SchemaField(
                    name="product_id", type="bigint", description="Product identifier"
                ),
                facet.SchemaField(
                    name="store_id", type="bigint", description="Store identifier"
                ),
                facet.SchemaField(
                    name="quantity", type="integer", description="Quantity of product sold"
                ),
                facet.SchemaField(
                    name="amount", type="bigdecimal", description="Amount of sale"
                ),
            ]
        ),
    },
)


source_managers = run.Dataset(
    namespace="retail_source",
    name="managers",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Manager identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of manager"
                ),
            ]
        ),
    },
)


source_stores = run.Dataset(
    namespace="retail_source",
    name="stores",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Store identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Store name"
                ),
                facet.SchemaField(
                    name="manager_id", type="bigint", description="Store manager"
                ),
            ]
        ),
    },
)


staged_brands = run.Dataset(
    namespace="retail_staged",
    name="brands",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Brand identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of brand"
                )
            ]
        ),
        "storage": facet.StorageDatasetFacet(
            storageLayer="staged", fileFormat="parquet"
        ),
        "symlinks": facet.SymlinksDatasetFacet(
            identifiers=[
                Identifiers(namespace="retail", name="brands", type="table")
            ]
        ),
        "ownership": facet.OwnershipDatasetFacet(
            owners=[
                Owners(name="chris@hyper.com", type="data analyst")
            ]
        ),
        "dataQualityAssertions": facet.DataQualityAssertionsDatasetFacet(
            assertions=[
                facet.Assertion(assertion="not_null", column="id", success=True)
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="brands",
                            field="id",
                        )
                    ]
                ),
                "name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="brands",
                            field="name",
                        )
                    ]
                ),
            }
        ),
    },
)

staged_products = run.Dataset(
    namespace="retail_staged",
    name="products",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Product identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of product"
                ),
                facet.SchemaField(
                    name="price", type="bigdecimal", description="Price of product"
                ),
                facet.SchemaField(
                    name="brand_id", type="varchar", description="Brand identifier"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="products",
                            field="id",
                        )
                    ]
                ),
                "name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="products",
                            field="name",
                        )
                    ]
                ),
                "price": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="products",
                            field="price",
                        )
                    ]
                ),
                "brand_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="products",
                            field="brand_id",
                        )
                    ]
                ),
            }
        )
    },
)


staged_managers = run.Dataset(
    namespace="retail_staged",
    name="managers",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Manager identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Name of manager"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="managers",
                            field="id",
                        )
                    ]
                ),
                "name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="managers",
                            field="name",
                        )
                    ]
                ),
            }
        ),
    },
)


staged_stores = run.Dataset(
    namespace="retail_staged",
    name="stores",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Store identifier"
                ),
                facet.SchemaField(
                    name="name", type="varchar", description="Store name"
                ),
                facet.SchemaField(
                    name="manager_id", type="bigint", description="Store manager"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="stores",
                            field="id",
                        )
                    ]
                ),
                "name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="stores",
                            field="name",
                        )
                    ]
                ),
            }
        ),
    },
)

staged_sales = run.Dataset(
    namespace="retail_staged",
    name="sales",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="id", type="bigint", description="Sale identifier"
                ),
                facet.SchemaField(
                    name="product_id", type="bigint", description="Product identifier"
                ),
                facet.SchemaField(
                    name="store_id", type="bigint", description="Store identifier"
                ),
                facet.SchemaField(
                    name="quantity", type="integer", description="Quantity of product sold"
                ),
                facet.SchemaField(
                    name="amount", type="bigdecimal", description="Amount of sale"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="sales",
                            field="id",
                        )
                    ]
                ),
                "product_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="sales",
                            field="product_id",
                        )
                    ]
                ),
                "store_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="sales",
                            field="store_id",
                        )
                    ]
                ),
                "quantity": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="sales",
                            field="quantity",
                        )
                    ]
                ),
                "amount": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_source",
                            name="sales",
                            field="amount",
                        )
                    ]
                ),
            }
        ),
    },
)


model_stores_dim = run.Dataset(
    namespace="retail_model",
    name="stores_dim",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="sk",
                    type="bigint",
                    description="Store dimension surrogate key",
                ),
                facet.SchemaField(
                    name="store_id", type="bigint", description="Store identifier"
                ),
                facet.SchemaField(
                    name="store_name", type="varchar", description="Store name"
                ),
                facet.SchemaField(
                    name="manager_id", type="bigint", description="Store manager id"
                ),
                facet.SchemaField(
                    name="manager_name", type="bigint", description="Store manager name"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "store_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="stores",
                            field="id",
                        )
                    ]
                ),
                "store_name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="stores",
                            field="name",
                        )
                    ]
                ),
                "manager_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="managers",
                            field="id",
                        )
                    ]
                ),
                "manager_name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="managers",
                            field="name",
                        )
                    ]
                ),
            }
        ),
    },
)

model_products_dim = run.Dataset(
    namespace="retail_model",
    name="products_dim",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="sk",
                    type="bigint",
                    description="Products dimension surrogate key",
                ),
                facet.SchemaField(
                    name="product_id", type="bigint", description="Product identifier"
                ),
                facet.SchemaField(
                    name="product_name", type="varchar", description="Product name"
                ),
                facet.SchemaField(
                    name="brand_id", type="bigint", description="Brand id"
                ),
                facet.SchemaField(
                    name="brand_name", type="bigint", description="Brand name"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "product_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="products",
                            field="id",
                        )
                    ]
                ),
                "product_name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="products",
                            field="name",
                        )
                    ]
                ),
                "brand_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="brands",
                            field="id",
                        )
                    ]
                ),
                "brand_name": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="brands",
                            field="name",
                        )
                    ]
                ),
            }
        ),
    },
)

model_sales_facts = run.Dataset(
    namespace="retail_model",
    name="sales_facts",
    facets={
        "schema": facet.SchemaDatasetFacet(
            fields=[
                facet.SchemaField(
                    name="sk",
                    type="bigint",
                    description="Surrogate key",
                ),
                facet.SchemaField(
                    name="sale_id", type="bigint", description="Sales id"
                ),
                facet.SchemaField(
                    name="quantity", type="int", description="Quantity sold"
                ),
                facet.SchemaField(
                    name="amount", type="bigdecimal", description="Amount sold"
                ),
                facet.SchemaField(
                    name="sale_time", type="datetime", description="Time when sale was made (UTC)"
                ),
                facet.SchemaField(
                    name="store_sk", type="bigint", description="Store dimension of the store where the sale happened"
                ),
                facet.SchemaField(
                    name="product_sk", type="bigint", description="Product dimension of the product sold"
                ),
            ]
        ),
        "columnLineage": facet.ColumnLineageDatasetFacet(
            {
                "sale_id": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="sales",
                            field="id",
                        )
                    ]
                ),
                "quantity": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="sales",
                            field="quantity",
                        )
                    ]
                ),
                "amount": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_staged",
                            name="sales",
                            field="amount",
                        )
                    ]
                ),
                "product_sk": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_model",
                            name="products_dim",
                            field="sk",
                        )
                    ]
                ),
                "store_sk": LineageFields(
                    transformationType="",
                    transformationDescription="",
                    inputFields=[
                        InputFields(
                            namespace="retail_model",
                            name="stores_dim",
                            field="sk",
                        )
                    ]
                ),
            }
        ),
    },
)

def run_load_source():
    c = client.OpenLineageClient.from_environment()
    run_id = str(uuid4())

    c.emit(
        run.RunEvent(
            eventType=run.RunState.START,
            eventTime="2021-11-03T10:53:52.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="load_source",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.RUNNING,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="load_source",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.COMPLETE,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="load_source",
            ),
            inputs=[source_managers, source_stores, source_brands, source_products, source_sales],
            outputs=[staged_managers, staged_stores, staged_brands, staged_products, staged_sales],
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

def run_build_dims():
    c = client.OpenLineageClient.from_environment()
    run_id = str(uuid4())

    c.emit(
        run.RunEvent(
            eventType=run.RunState.START,
            eventTime="2021-11-03T10:53:52.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_dims",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.RUNNING,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_dims",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.COMPLETE,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_dims",
            ),
            inputs=[staged_managers, staged_stores, staged_brands, staged_products],
            outputs=[model_stores_dim, model_products_dim],
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )


def run_build_facts():
    c = client.OpenLineageClient.from_environment()
    run_id = str(uuid4())

    c.emit(
        run.RunEvent(
            eventType=run.RunState.START,
            eventTime="2021-11-03T10:53:52.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_facts",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.RUNNING,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_facts",
            ),
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

    c.emit(
        run.RunEvent(
            eventType=run.RunState.COMPLETE,
            eventTime="2021-11-03T10:53:53.427Z",
            run=run.Run(runId=run_id),
            job=run.Job(
                namespace="retail",
                name="build_facts",
            ),
            inputs=[staged_sales, model_stores_dim, model_products_dim],
            outputs=[model_sales_facts],
            producer="https://github.com/OpenLineage/OpenLineage/tree/0.0.1/client/python",
        )
    )

if __name__ == '__main__':
    run_load_source()
    run_build_dims()
    run_build_facts()
