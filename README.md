# Authboss Sample
A sample implementer of authboss.

This is a simple blogging engine with a few basic features:

- Authentication provided by Authboss (all modules enabled, despite conflict between remember & expire)
- Overridden (pretified) Authboss views.
- CRUD for an in-memory storage of blogs.
- Flash Messages
- XSRF Protection

**Disclaimer:** This sample is NOT a seed project. Do not use it as one. It is used as an example of how to use the Authboss API.
This means if you copy-paste code from this sample you are likely opening yourself up to various security holes, bad practice,
and bad design. It's a demonstration of the surface API of Authboss and how the library can be used to make a functioning web
project, to use this sample as anything else is malpractice.

## Get started

1. Copy env_sample.sh to .env and edit OAuth client id/secret

2. Build app
```bash
go get gopkg.in/authboss.v1
git clone https://github.com/go-authboss/authboss-sample
cd authboss-sample
go build
```

3. Run app

```bash
source .env
./authboss-sample
```

By default it should work by opening `http://localhost:3000/`

Note that email and some sample users will be logged in console