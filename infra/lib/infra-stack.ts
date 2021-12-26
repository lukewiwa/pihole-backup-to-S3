import { Stack, StackProps, CfnParameter, CfnOutput, Tags, RemovalPolicy } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as iam from "aws-cdk-lib/aws-iam"

export class InfraStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const PROJECT_TAG: [string, string] = ["project", "piholebackup"]

    const uploadBucketName = new CfnParameter(this, "piholeBackupBucketName", {
      type: "String",
      description: "The name of the Amazon S3 bucket where uploaded files will be stored."});

    const bucket = new s3.Bucket(this, "PiholeBackup", {
      bucketName: uploadBucketName.valueAsString,
      removalPolicy: RemovalPolicy.DESTROY
    })
    Tags.of(bucket).add(...PROJECT_TAG)

    const user = new iam.User(this, "PiholeBackupUser")
    Tags.of(user).add(...PROJECT_TAG)

    bucket.grantRead(user)
    bucket.grantPut(user)

    const accessKey = new iam.CfnAccessKey(this, "PiholeBackupAccessKey", {userName: user.userName})

    new CfnOutput(this, 'accessKeyId', { value: accessKey.ref });
    new CfnOutput(this, 'secretAccessKey', { value: accessKey.attrSecretAccessKey });
  }
}
