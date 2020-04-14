
This is the first blog post from a series of posts where I'm going to show the options I tried to build my site.   
Hopefully, this will give you some hints on what it takes to build a simple site/blog using different approaches.   
* Although my initial intention was to build the site in GO, I wanted to give a try to this wonderfull tool, Hugo,  which is written in GO. If you are looking to build a static site this can be a satisfactory and sufficient solution. On the next posts I'll focus more on what GO standard library offer to us and the available frameworks to build the website.

### Build your site

The first thing I tried was to generate a static site using the Hugo framework.   
<b>What is Hugo?</b>   
Based on their logo "<b> Hugo is one of the most popular open-source static site generators. With its amazing speed and flexibility, Hugo makes building websites fun again. </b>"   
Hugo is written in Go with support for multiple platforms, but you don’t need to install Go to enjoy Hugo.

The best part is that you don't need to know any programming languange in order to build your own site.
All files are generated by Hugo and copied in the public folder. Once done, the site is ready to be hosted on your local or cloud server. 

<center>
    <a href="/images/hugo.jpg" target="_blank"><img src="/images/hugo.jpg" /></a>
</center>

<b>Install Hugo</b>

Hugo installation for different platforms is very well covered at: https://gohugo.io/getting-started/installing

I'm using Linux, and I opted for snap package.

```bash
snap install hugo --channel=extended
```

<b> Create the site </b>

```bash
hugo new site site_name
```

Add a theme, I tried several but I liked more the [Future Imperfect Slim] which is a modern and clean theme:
[Future Imperfect Slim]: https://themes.gohugo.io/hugo-future-imperfect-slim/

```bash
cd themes/
git clone https://github.com/pacollins/hugo-future-imperfect-slim.git
cd ..
```

Most of the themes have an exampleSite, which you can use as an example when you create your own site. More than that, the template can be tried directly on the Hugo site template presentation.

Copy the configuration file and tweak it for your needs

```bash
cp themes/hugo-future-imperfect-slim/exampleSite/config.toml .
```

Create your content:

```bash
hugo new about/_index.md
hugo new blog/examples.md
```

Populate your markdown files with the site content. If you are not familiar with editing markdown files,this [Markdown Cheatsheet] can help you get off the ground.
[Markdown Cheatsheet]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet

The changes can be seen in real time, Hugo will build your site and host a server locally. You can view this live at localhost:1313.

```bash
hugo server
```

At the end, the hugo command will build your site by generating the html files and copy static content in the public folder.

```bash
hugo
```

### Host your site on Cloud (GCP) 

In my case I opted for Cloud Storage from Google Cloud platform. As part of the Google Cloud Platform Free Tier, Cloud Storage provides resources that are free to use up to specific limits.   
As a prerequisite you need a valid GCP account and use the Web Console or install Cloud SDK on your platform of choice to manage your platform.    
The configuration and instalation of Cloud SDK is out of the scope of this tutorial,
I found the [Quickstarts] very helpful in setting up the environtment.
[Quickstarts]: https://cloud.google.com/sdk/docs/quickstarts    
Moving forward I assume that you had minimal exposure to Cloud platforms and have everything setted up.

All the steps I followed below can be performed directly from web Console, you can check the official tutorial, [here].
[here]: https://cloud.google.com/storage/docs/hosting-static-website

* Prerequisite, to create a bucket that uses a domain name, you must establish that you are authorized to use the domain name.  For this, follow the steps shown in [this page].
[this page]: https://cloud.google.com/storage/docs/domain-name-verification#verification 

<b> Create your storage bucket where the site will be hosted </b>

```bash
gsutil mb gs://www.your_site_name.com
```

* replace www.your_site_name.com with your real site name

<b> Upload your site's files </b>   
In our case, hugo generated public folder should be uploaded to the bucket. 

```bash
gsutil -m rsync -R public gs://www.your_site_name.com
```

<b> Share your files </b>   
In order to allow everyone to acces your site you have to make the objects within the bucket publicly readable

```bash
gsutil iam ch allUsers:objectViewer gs://www.your_site_name.com
```

<b> Define the site entry point and error page </b>
You can assign an index page suffix, which is controlled by the MainPageSuffix property and a custom error page, which is controlled by the NotFoundPage property. Assigning either is optional, but without an index page, nothing is served when users access your top-level site, for example, http://www.your_site_name.com.

```bash
gsutil web set -m index.html -e 404.html gs://www.your_site_name.com
```

### Change the DNS record

Last step in this tutorial is to create the CNAME record. A CNAME record is a type of DNS record. It directs traffic that requests a URL from your domain to the resources you want to serve, in this case objects in your Cloud Storage buckets.   
In my case I'm using GoDaddy to manage my DNS records.   
I have added a Type: CNAME , Name: www, Value: c.storage.googleapis.com and saved.   
It may take awhile until the dns will propagate.   
