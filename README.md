# Local Content Share

For a better explanation, check out my [blog post](https://blog.tanishq.page/posts/homelab-local-dumpster/) for deploying and using this app in a home lab.

This application can be used for &rarr;

- keep a record of pasted markdown text available for copying by any machine in the local network (also renderable)
- keep a record of uploaded files available for download by any machine in the local network
- keep a record of links available for viewing/visiting by any machine in the local network
- render a given markdown text in github-flavor for printing/viewing

Assuming that the application is deployed at a system with the hostname `galaxy`, visit `http://galaxy.local` for using the app.

For running the docker image, do the following &rarr;

```bash
docker run --name local_dumpster --rm -p 80:5000 -d -t tanq16/local_dumpster:main
# use tage :main_arm for ARM64 image (useful for apple silicon and raspberry pi)
```

This will launch the container as a daemon and it'll be reachable at port 80 on the host machine.
