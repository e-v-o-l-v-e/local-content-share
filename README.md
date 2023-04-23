# MD Share & Render

### Usage

This application can be used for 2 things &rarr;

- keep a record of pasted markdown text available for copying by any machine in the local network
- render the markdown text in github-flavor for printing/viewing

Assuming that the application is deployed at a system with the hostname `galaxy`, the following actions can be taken &rarr;

- Visit `http://galaxy.local` and add markdown text. Entries will be persisted and listed at the bottom with the first line visible and options to view raw, github-flavored render, github-flavored render with dark theme, and delete.
- Visit `http://galaxy.local/print` to simply paste and render markdown text in github-flavor or visit `http://galaxy.local/print_dark` to render in github-flavor dark mode. 

### Installation

For building and running the docker image, do the following &rarr;

```bash
git clone https://github.com/tanq16/share-n-render
cd share-n-render
docker build -t local_dumpster .
docker run --name md_dumpster --rm -p 80:5000 -d -t local_dumpster
```

This will launch it as a daemon container and it would be reachable at port 80 on the host machine.

The recommended method of installation and deployment is via docker. If local installation is needed, refer to the modules installed within the Dockerfile, and install them within a virtual environment before launching the flask app.
