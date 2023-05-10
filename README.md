<div align="center">
  <img src="./hizuru-overview.jpg" alt="hizuru icon" width="200px" />
</div>

# hizuru
Simple Windows Terminal background image changer.  
With a few settings and a single command, you can switch to your favorite background image at any time.  

## Demo
https://github.com/yudaishimanaka/hizuru/assets/11958380/6dec801a-205e-4dbd-9a25-9bac4f15ac6f

## Requirement
All you need is a [Windows Terminal](https://github.com/microsoft/terminal).

## Install
1. Download from the [release](https://github.com/yudaishimanaka/hizuru/releases/tag/v1.0.0) according to your computer architecture.
2. Unzip and register the PATH of the command.
3. Set the `HIZURU_IMAGE_PATH` environment variable.
    ```powershell
    [Environment]::SetENvironmentVariable("HIZURU_IMAGE_PATH", "C:\hoge\fuga", "User")
    ```

## Usage
Place the background image directly under the `HIZURU_IMAGE_PATH` directory, then simply run the command.  

Only one subcommand and flag.  

Change background image.  
`hizuru change`

If you are using the preview version of the Windows Terminal.  
`hizuru --preview change`

## Licence
[MIT](https://github.com/yudaishimanaka/hizuru/blob/master/LICENSE)

## Author
[yudaishimanaka](https://github.com/yudaishimanaka)
