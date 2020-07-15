# Tidy

Tidy is a desktop application created to manage the deletion of all those local branches that have been already merged into a target branch.

This is a little tool i built in order to keep my local repository clean, since the current workflow in my job is to create a branch per feature and afterwards merge the aforementioned branch into a development branch, we end up having lots of branches laying around in the local repository which is not good.

<p align="center">
  <img src="./assets/tidy.png" alt="tidy" />  
</p>

<table border="0" cellspacing="0" cellpadding="0" style="border-collapse: collapse; border: none;">
  <tr>
    <td><img alt="GitHub" src="https://img.shields.io/github/license/wwleak/tidy?style=for-the-badge"></td>
    <td><img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/wwleak/tidy?style=for-the-badge"></td>
    <td><img alt="GitHub top language" src="https://img.shields.io/github/languages/top/wwleak/tidy?style=for-the-badge"></td>
    <td><img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/wwleak/tidy?style=for-the-badge"></td>
    <td><img alt="GitHub stars" src="https://img.shields.io/github/stars/wwleak/tidy?style=for-the-badge"></td>
  </tr>
</table>

Gtkinspector is a tool which allow us to inspect all aspects of gtk windows, it's really useful for desktop apps development. It helped me to add custom styling to the gtk windows just with plane css.

<p align="center">
  <img src="assets/inspector.png" name="inspector" />
</p>

## :rocket: How to run it ?

Clone the repository

```bash
 git clone https://github.com/wwleak/tidy.git
```
Build the binary

```bash
  go build -o tidy
```

All you need to do is to add some configuration, in the settings section you have to provide where in the filesystem is laying your repository along with the target branch you wish to know your branches were merged to.

<p align="center">
  <img src="assets/settings.png" name="inspector" />
</p>

That's it.

Another option less user friendly is to let a bash script do the job, you can find a script in my gists
called [**`cleanup.sh`**](https://gist.github.com/wwleak/0f6238417e64512f81c1f40759115b3e)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## :pushpin: License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

<p align="right">MADE WITH ❤ BY ABRAHAM</p>
