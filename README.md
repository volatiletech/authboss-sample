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

# Getting Started

To get started, clone the repo and execute the following commands

```
cd authboss-sample
export GOPATH=`pwd`
go get
```

After all the packages have been installed (you should see a `bin`, `pkg`, and `src` directories), execute:

```
bin/authboss-sample
```

Finally, open `http://127.0.0.1:3000/auth/login` in your browser and login using:

```
username: zeratul@heroes.com
password: 1234
```

You should see the blog posts.
