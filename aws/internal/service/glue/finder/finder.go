package finder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	tfglue "github.com/terraform-providers/terraform-provider-aws/aws/internal/service/glue"
)

// RegistryByID returns the Registry corresponding to the specified ID.
func RegistryByID(conn *glue.Glue, id string) (*glue.GetRegistryOutput, error) {
	input := &glue.GetRegistryInput{
		RegistryId: tfglue.CreateAwsGlueRegistryID(id),
	}

	output, err := conn.GetRegistry(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SchemaByID returns the Schema corresponding to the specified ID.
func SchemaByID(conn *glue.Glue, id string) (*glue.GetSchemaOutput, error) {
	input := &glue.GetSchemaInput{
		SchemaId: tfglue.CreateAwsGlueSchemaID(id),
	}

	output, err := conn.GetSchema(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SchemaVersionByID returns the Schema corresponding to the specified ID.
func SchemaVersionByID(conn *glue.Glue, id string) (*glue.GetSchemaVersionOutput, error) {
	input := &glue.GetSchemaVersionInput{
		SchemaId: tfglue.CreateAwsGlueSchemaID(id),
		SchemaVersionNumber: &glue.SchemaVersionNumber{
			LatestVersion: aws.Bool(true),
		},
	}

	output, err := conn.GetSchemaVersion(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PartitionByValues returns the Partition corresponding to the specified Partition Values.
func PartitionByValues(conn *glue.Glue, id string) (*glue.Partition, error) {

	catalogID, dbName, tableName, values, err := tfglue.ReadAwsGluePartitionID(id)
	if err != nil {
		return nil, err
	}

	input := &glue.GetPartitionInput{
		CatalogId:       aws.String(catalogID),
		DatabaseName:    aws.String(dbName),
		TableName:       aws.String(tableName),
		PartitionValues: aws.StringSlice(values),
	}

	output, err := conn.GetPartition(input)
	if err != nil {
		return nil, err
	}

	if output == nil && output.Partition == nil {
		return nil, nil
	}

	return output.Partition, nil
}
