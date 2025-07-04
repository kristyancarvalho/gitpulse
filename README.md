<p align="center">
  <h1 align="center">GitPulse</h1>

  <p align="center">
    <img src="https://img.shields.io/github/repo-size/kristyancarvalho/gitpulse?color=success" alt="Repo size" />
    <img src="https://img.shields.io/github/last-commit/kristyancarvalho/gitpulse" alt="Last commit" />
    <img src="https://img.shields.io/github/stars/kristyancarvalho/gitpulse?style=social" alt="Stars" />
    <img src="https://img.shields.io/github/forks/kristyancarvalho/gitpulse?style=social" alt="Forks" />
  </p>
</p>

---

Webservice built in Go that generates an image displaying the **latest public GitHub contribution** of a user, showing the repository name and its main language.

## Example

```html
<img src="https://gitpulse.kristyan.dev/api/last-project?username=kristyancarvalho&color=%23b600ff" alt="Latest contribution"/>
````

**Result:** <br /> <img src="https://gitpulse.kristyan.dev/api/last-project?username=kristyancarvalho&color=%23b600ff" alt="Latest contribution" />

---

## Parameters

| Parameter  | Required | Description                                                    |
| ---------- | -------- | -------------------------------------------------------------- |
| `username` | ✅ Yes    | GitHub username                                                |
| `color`    | ❌ No     | Hexadecimal color (e.g., `%23ff9900`) for the image text color |

---

## Example usage in a GitHub README

```markdown
![Latest contribution](https://gitpulse.kristyan.dev/api/last-project?username=your_username)
```

---

## Built with

* [Golang](https://golang.org/)
* [GitHub REST API v3](https://docs.github.com/en/rest)
* [Vercel](https://vercel.com/) (deployment)
