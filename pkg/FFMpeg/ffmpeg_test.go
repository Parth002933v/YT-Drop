package FFMpeg

import (
	"fmt"
	"testing"
)

func TestExtractChapters2(t *testing.T) {
	text := `⭐️ Contents ⭐️
(0:03:51) Introduction
(0:03:51) Prerequisites
(0:04:21) Setting up our project and overview
(0:07:03) Root route explained and linking our CSS
(0:08:22) Creating your first route and render via outlet
(0:10:36) Creating Dynamic Routes in Remix
(0:14:12) Setting up contact detail page
(0:15:08) Using the loader function to load data
(0:20:02) Loading single-user based on id via params
(0:24:48) Setting up Strapi, a headless CMS
(0:27:56) Strapi Admin overview and creating our first collection type
(0:33:20) Fetching all contacts from our Strapi endpoint
(0:38:17) Fetching single contact
(0:40:17) Adding the ability to add a new contact
(0:54:41) Form validation with Zod and Remix
(1:03:31) Adding the ability to update contact information
(1:16:25) Programmatic navigation using useNavigate hook
(1:18:15) Implementing the delete contact functionality
(1:25:53) Showing a fallback page when no items are selected
(1:27:25) Handling errors in Remix with error boundaries
(1:34:04) Implementing mark contact as a favorite
(1:37:33) Implementing search with Remix and Strapi filtering
(1:58:51) Submitting our form programmatically on input change
(2:00:39) Implementing loading state via useNavigation hook
(2:05:45) Project review and some improvement
(2:06:55) Styling active link in Remix
(2:09:17) Using useFetcher hook for form submission 
(2:11:08) Throwing errors in Remix
(2:15:41) Closing thought and where to find hel`

	videoDuration := int64(301043230) // Example video duration in milliseconds

	chapters, err := ExtractChapters(text, videoDuration)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	// Print chapters to verify output
	for _, chapter := range chapters {
		fmt.Printf("Start: %dms, End: %dms, Title: %s\n", chapter.StartTime, chapter.EndTime, chapter.Title)
	}
}

func TestExtractChapters3(t *testing.T) {
	test2 := `⭐️ Contents ⭐️
#1   (00:00:00) CSS tutorial for beginners 🎨
#2   (00:11:00) fonts 🆒
#3   (00:14:20) borders 🔲
#4   (00:16:56) background 🌆
#5   (00:20:52) margins 📏
#6   (00:25:44) float 🎈
#7   (00:29:01) position 🎯
#8   (00:34:58) pseudo classes 👨‍👧‍👦
#9   (00:40:47) shadows 👥
#10 (00:43:43) icons 🏠
#11 (00:46:45) transform 🔄
#12 (00:50:54) animation 🎞️`
	videoDuration := int64(301043230) // Example video duration in milliseconds

	chapters, err := ExtractChapters(test2, videoDuration)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	// Print chapters to verify output
	for _, chapter := range chapters {
		fmt.Printf("Start: %dms, End: %dms, Title: %s\n", chapter.StartTime, chapter.EndTime, chapter.Title)
	}
}

func TestExtractChaptersInError(t *testing.T) {
	text := `✏️ Course developed by Andrew Brown of ExamPro.  ‪@ExamProChannel‬ 
⭐️ Contents ⭐️
0:00:00 Introduction
0:34:47 Setup
0:52:38 Amazon S3
10:52:02 AWS API
12:19:52 VPC
0:34:47 Setup22222222222222
17:33:42 IAM
19:14:03 EC2
21:13:27 AMIs
21:37:10 ASG
21:50:14 ELB
21:57:20 Route53
22:19:29 AWS Global Accelerator
22:21:00 CloudFront
22:30:24 EBS
22:45:34 EFS
22:50:38 FSx
22:54:24 AWS Backup
22:56:29 AWS Snow Family
23:07:07 AWS Transfer Family
23:09:31 AWS Migration Hub
23:15:35 AWS Data Sync
23:24:17 DMS
23:59:42 AWS Auto Scaling
24:16:59 AWS Amplify
24:37:15 Amazon AppFlow
24:53:39 AppSync
25:18:48 AWS Batch
25:46:37 OpenSearch Service
26:09:43 Device Farm
26:22:11 QLDB
26:24:01 Elastic Transcoder
26:52:21 AWS MediaConvert
27:02:09 SNS
27:43:05 SQS
28:44:00 Amazon MQ
29:32:34 Service Catalog
29:40:04 CloudWatch and EventBridge
30:16:36 Lambda
31:49:51 AWS Step Functions
32:48:57 AWS Compute Optimizer
32:59:19 Elastic Beanstalk
34:32:38 Kinesis
34:59:52 ElastiCache
35:51:13 MemoryDB
36:21:52 CloudTrail
37:19:23 Redshift
37:37:50 Athena
37:53:46 ML Managed Services
40:43:04 AWS Data Exchange
40:47:11 AWS Glue
41:27:04 Lake Formation
41:29:41 API Gateway
41:44:09 RDS
42:56:19 Aurora
19:33:29 DocumentDB
44:29:11 DynamoDB
21:10:04 Amazon Keyspaces
45:17:30 Neptune
45:35:00 ECR
45:39:18 ECS
46:02:27 EKS Cloud
46:21:45 KMS
46:32:00 AWS Audit Manager
46:40:23 ACM
46:58:57 Cognito
47:08:33 Amazon Detective
47:16:42 AWS Directory Service
47:22:47 AWS Firewall Manager
47:29:18 AWS Inspector
47:39:57 Amazon Macie
47:49:00 AWS Security Hub
47:53:37 AWS Secrets Manager
48:35:40 AI Dev Tools
48:59:17 Amazon MSK
49:29:32 AWS Shield
49:33:29 AWS WAF
49:37:48 CloudHSM
49:41:59 AWS Guard Duty
49:46:10 Health Dashboards
49:47:42 AWS Artifact
49:50:33 Storage Gateway
50:10:55 EC2 Pricing Models`
	videoDuration := int64(301043230) // Example video duration in milliseconds

	chapters, err := ExtractChapters(text, videoDuration)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	// Print chapters to verify output
	for _, chapter := range chapters {
		fmt.Printf("Start: %dms, End: %dms, Title: %s\n", chapter.StartTime, chapter.EndTime, chapter.Title)
	}
}
