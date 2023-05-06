# MD Share & Render

For a better explanation, check out my [blog post](https://blog.tanishq.page/posts/homelab-md-dumpster/) for deploying and using this app in a home lab.

This application can be used for &rarr;

- keep a record of pasted markdown text available for copying by any machine in the local network (also renderable)
- keep a record of uploaded files available for download by any machine in the local network
- keep a record of links available for viewing/visiting by any machine in the local network
- render a given markdown text in github-flavor for printing/viewing

Assuming that the application is deployed at a system with the hostname `galaxy`, visit `http://galaxy.local` for using the app.

For building and running the docker image, do the following &rarr;

```bash
git clone https://github.com/tanq16/share-n-render
cd share-n-render
docker build -t local_dumpster .
docker run --name md_dumpster --rm -p 80:5000 -d -t local_dumpster
```

This will launch it as a daemon container and it would be reachable at port 80 on the host machine.
