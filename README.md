<!-- markdownlint-configure-file { "MD004": { "style": "consistent" } } -->
<!-- markdownlint-disable MD033 -->

#

<p align="center">
  <picture>
    <img src="https://equalsgibson.github.io/concur/concur.jpg" width="410" height="205" alt="Concur">
  </picture>
    <br>
    <strong>Easily run concurrent goroutines safely</strong>
</p>

<!-- markdownlint-enable MD033 -->

- **Easy to use**: Get up and running with the library in minutes
- **Intuitive**: Get access to a powerful pattern within software development, without encountering common pitfalls or mistakes!
- **Actively developed**: Ideas and contributions welcomed!

---

<div align="right">

[![Go][golang]][golang-url]
[![Go Reference][goref]][goref-url]
[![Go Report Card][goreport]][goreport-url]

</div>

## Getting Started

`concur` is a Go library that allows you to safely and easily implement a concurrency pattern in your codebase. Install the latest version, and get up and running in minutes.

This library is currently used by the Five9-GO repository, to provide a reliable, safe way to interact with the Five9 API Websocket service.

### What is concurrency, and why is it useful?

```
Concurrency is the composition of independently executing computations.

Concurrency is a way to structure software, particularly as a way to write clean code that interacts well with the real world.

It is not parallelism.

	- Rob Pike, 2012
```

Simply put, concurrency in software allows you to create **fast**, **robust** systems that can be relied upon to be **consistent**. Some examples of what concurrency can do, as provided by the Go team:

- [A Prime Number Sieve](https://go.dev/play/p/9U22NfrXeq)
- [An RSS Feed fetcher](https://cs.opensource.google/go/x/website/+/master:_content/talks/2013/advconc/realmain/realmain.go)

### Common pitfalls of concurrency that this library prevents

- Race conditions
- Deadlocks
- Unpredictable or "flaky" tests when testing concurrent data models (such as Websockets, which are inherently asynchronous messages being sent back and forth between systems)

### Install

```shell
go get github.com/equalsgibson/concur@v0.0.3
```

### Example of ASyncReader

Run the example ASyncReader with a mock asynchronous resource, along with a CPU and Memory profile.

https://github.com/equalsgibson/concur/blob/e3d8905e740043133e92e11918fd1c2e47d7e0e3/examples/async_reader_timer.go/main.go#L1-L109

Expected Output:

```bash
cgibson@wsl-ubuntuNexus:~/concur$ go run examples/async_reader_timer.go/main.go -cpuprofile cpu.prof -memprofile mem.prof
2024/10/29 00:01:00 Current Iteration: 1 - Current goroutines: 4
2024/10/29 00:01:01 Current Iteration: 2 - Current goroutines: 4
...
2024/10/29 00:01:03 Current Iteration: 299 - Current goroutines: 4
2024/10/29 00:01:03 Current Iteration: 300 - Current goroutines: 4
2024/10/29 00:01:05 Got an error response from Loop fetch function: end of the example - thanks for using the concur package
```

You can then analyze the profiles created, using the [pprof](https://pkg.go.dev/runtime/pprof#hdr-Profiling_a_Go_program) tool

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, get inspired, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<!-- CONTACT -->

## Contact

[Chris Gibson (@equalsgibson)](https://github.com/equalsgibson)

Project Link: [https://github.com/equalsgibson/concur](https://github.com/equalsgibson/concur)

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[golang]: https://img.shields.io/badge/v1.21-000?logo=go&logoColor=fff&labelColor=444&color=%2300ADD8
[golang-url]: https://go.dev/
[goref]: https://pkg.go.dev/badge/github.com/equalsgibson/concur.svg
[goref-url]: https://pkg.go.dev/github.com/equalsgibson/concur
[goreport]: https://goreportcard.com/badge/github.com/equalsgibson/concur
[goreport-url]: https://goreportcard.com/report/github.com/equalsgibson/concur
