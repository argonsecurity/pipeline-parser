package models

type BuildPipelines struct {
	Default      []*Step  `yaml:"default"`                 // The default pipeline runs on every push to the repository, unless a branch-specific; pipeline is defined.; You can define a branch pipeline in the branches section.; ; Note: The default pipeline doesn't run on tags or bookmarks.
	Branches     *StepMap `yaml:"branches,omitempty"`      // Defines a section for all branch-specific build pipelines. The names or expressions in; this section are matched against:; ; * branches in your Git repository; * named branches in your Mercurial repository; ; You can use glob patterns for handling the branch names.
	Tags         *StepMap `yaml:"tags,omitempty"`          // Defines all tag-specific build pipelines.; ; The names or expressions in this section are matched against tags and annotated tags in; your Git repository.; ; You can use glob patterns for handling the tag names.
	Bookmarks    *StepMap `yaml:"bookmarks,omitempty"`     // Defines all bookmark-specific build pipelines.; ; The names or expressions in this section are matched against bookmarks in your Mercurial; repository.; ; You can use glob patterns for handling the tag names.
	PullRequests *StepMap `yaml:"pull-requests,omitempty"` // A special pipeline which only runs on pull requests. Pull-requests has the same level of; indentation as branches.; ; This type of pipeline runs a little differently to other pipelines. When it's triggered,; we'll merge the destination branch into your working branch before it runs. If the merge; fails we will stop the pipeline.
	Custom       *StepMap `yaml:"custom,omitempty"`        // Defines pipelines that can only be triggered manually or scheduled from the Bitbucket; Cloud interface.
}
