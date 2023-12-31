
# Change Log
## [Unreleased] - 2024-09-28

### Added
- CHANGELOG.md to track changes

### Changed
- Add error handling for ResponseWriter failures to backend.
- Make `main` the default branch

### Fixed
- Updated search function to use a regular expression and disable case sensitivity. I'm using the `FindAllIndex` method for the `SuffixArray` of the complete works to allow the use of a regular expression. Since it returns start and end indexes, I want to ensure the results are still the same, so I'm only using the start index.
  I create the regex using MustCompile for brevity, I need to investigate if this could cause issues of not-nicely formatted strings.
- Limited the search results to the first 20.
- Make `Load More` button add more results. I decided to load all existing results again as well which might become a performance problem in the future and will have to be adapted then.