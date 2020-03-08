package blog

type BlogActor interface {
	ActorTest() (bool, string)
	PublishBlog(publishTime, title, content string) (bool, string)
	UpdateBlog(oldtitle, firstPublishTime, newTitle, content string) (bool, string)
	DeleteBlog(blogid string) (bool, string)
}
