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
