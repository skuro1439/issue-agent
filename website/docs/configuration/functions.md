# Functions
```
get_pull_request: Get a GitHub Pull Request
    pr_number
        Pull Request Number to get

search_files: Search for files containing specific keyword (e.g., "xxx") within a directory path recursively
    keyword
        The keyword to search for.
    path
        The path to search within its directory

list_files: List the files within the directory like Unix ls command. Each line contains the file mode, byte size, and name
    path
        The valid path to list within its directory

put_file: Put new content to the file
    output_path
        Path of the file to be changed to the new content
    content_text
        The new content of the file

get_web_page_from_url: Get the web page from the URL
    url
        The URL to get the Web page. More than 80000 characters are cut off

get_web_search_result: Get a list of results from an Internet search conducted with keywords. You should get the page information from the url of the result next.
    keyword
        Keyword to search for on the Internet

open_file: Open the file full content
    path
        The path of the file to open

modify_file: Modify the file at output_path with the contents of content_text. Modified file must be full content including modified content
    output_path
        Path of the file to be modified to the new content
    content_text
        The new content of the file

submit_files: Submit the modified files by GitHub Pull Request
    commit_message_short
        Short Commit message indicating purpose to change the file
    commit_message_detail
        Detail commit message indicating changes to the file
    pull_request_content
        Pull Request Content

```
