
# Table of contents

- [introduction](#introduction)
  - [overview](#overview)
  - [getting started](#getting-started)



# Introduction

## overview
wuxia is a potable static web solution. Wuxia allows you to create, build
deploy and scale static websites.

The main motivation behind this is to open up ways for people in my country to
write technical documents  i.e books, articles, etc.

Wuxia is made up of components that when combined together forms a powerful
solution for the static web fans. It ships with a static site generator, static
site server, a `git push` based static site deployment , a simple git repository
host service( support multi user and multi repositories ) and a multi domain static
site host all in a single binary.

I built this for personal use, but if you find it useful please do the
followwing, [star the project](https://github.com/gernest/wuxia) , [fork the project](https://github.com/gernest/wuxia), [follow the author](https://github.com/gernest)
and in case you encounter any bugs then [open an issue](https://github.com/gernest/wuxia/issues).

# Generator
Wuxia ships with a static site generator. The implementation of generator is as
a package which means it can be used independently  with another project.

### Why another static generator?
This question is inevitable, there is already a ton load of static site
generators out there, so why did I decide we need a room for one more?

- Seamless integration with the current system. Since the system is written in
  Go, It was obvious using a Go based static site generator would be a good
  idea except for the fact that most static site generators are applications and
  not libraries. Wuxia generator is a library that you can integrate into a
  bigger system.

- Metrics: I needed a way not only to have benchmarks  but all the metrics that
  one can gather about the static content generation process.

- Ease to extend with plugins: I wanted a generator that can be extended in
  functionality and not a black box. This throws many Go options out of the
  window. The current implementation uses javascript as a scripting language.
  This opens up a wide array of good stuffs to the static web generation
  experience.


## Overview of features of the generator

- Convention plus configuration. It comes with sane defaults, sane conventions
  of the way the project should be structured. On to p of this, it gives full
  freedom if the user is advanced enough to configure it in way that it suits
  his/her needs.

- Themes : Support customizing of project's look and feel. The default template
  language is is Go. Effort is underway to add popular javascript templating
  engines.

- Javascript Plugins: Write plugins in javascript. You can tap almost every
  aspect of the generation process through javascript. The generator comes with
  a simple module loader `require` which you can use to break down your modules
  into separate files.

- Metrics: The generator is wired with exhaustive metrics. You can have full
  understanding of the resources used and perfomrance of rendering your project.


## The structure of a wuxia static website project

```
├── scripts
│   ├── init
│   │   └── index.js
│   ├── plan
│   │   └── index.js
│   └── plugin
│       └── index.js
├── src
│   ├── front_matter.md
│   └── home.md
├── templates
│   └── index.html
├── LICENCE
├── README.md
└── config.json
```
