# parser-web-crawler
This is a Coding challenge to show up my skill as a Golang Senior

# Challenge Description:
We’d like you to write a simple web crawler in Golang.
Given a starting URL, the crawler should visit each URL it finds on the same domain. It should
print each URL visited, and a list of links found on that page. The crawler should be limited to
one subdomain – so when you start with https://parserdigital.com/, do not follow external
links, for example to facebook.com or community.parserdigital.com.
We would like to see your own implementation of a web crawler. Please do not use frameworks
like scrappy or go-colly which handle all the crawling behind the scenes or someone else’s
code. You are welcome to use libraries to handle things like HTML parsing.
Ideally, write it as you would a production piece of code. This exercise is not meant to show us
whether you can write code – we are more interested in how you design software. This means
that we care less about a fancy UI or sitemap format, and more about how your program is
structured: the trade-offs you’ve made, what behaviour the program exhibits, and your use of
concurrency, test coverage, and so on.
Once you have submitted your task, we will then schedule a session with an engineer, during
which we all will discuss your implementation.
When you’re ready, please submit your solution as a ZIP file.


## How to Build/Run it
go build

## We can use the default configuration and run it by:
./parser-web-crawler 

### Custom the Web Crawler by change the parameters:
-n for the number of workers

-url for the URL to start crawling


#### Knowing how to custom it, we can do:


./parser-web-crawler -n 35 -url https://cesarcedanov.com/

./parser-web-crawler -url https://cesarcedanov.com/

./parser-web-crawler -n 35






