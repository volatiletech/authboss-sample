# Authboss Sample
A sample implementation of authboss.

This is a simple blogging engine with a few basic features:

- Authentication provided by Authboss (all modules enabled with the exception of expire)
- Some examples of overridden Authboss views.
- CRUD for an in-memory storage of blogs.
- Flash Messages
- CSRF Protection (including authboss routes)
- Support for API style JSON requests and responses (-api flag)
- Various levels of debugging to see what's going wrong (-debug* flags)

Uses the following default libraries:

- https://github.com/volatiletech/authboss-renderer
- https://github.com/volatiletech/authboss-clientstate

# Disclaimer

This sample is **NOT** a seed project. Do not use it as one.
It is used as an example of how to use the Authboss API. This means if
you copy-paste code from this sample you are likely opening yourself
up to various security holes, bad practice, and bad design.
It's a demonstration of the surface API of Authboss and how the library
can be used to make a functioning web project.
