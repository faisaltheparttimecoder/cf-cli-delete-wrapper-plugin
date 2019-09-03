# cf-cli-delete-wrapper-plugin

A cf cli plugin that is a wrapper on top of "cf delete" command.
This plugin wrapper provides two commands that helps to easy cleanup of apps.
    
   1. You can now delete multiple apps via single command
   2. Command that auto detects the app name from manifest and deletes them
    
# Installation

You can use either of the below method to install the plugin.

**Option 1:** 

   1. Download the latest plugin from the [release](https://github.com/faisaltheparttimecoder/cf-cli-delete-wrapper-plugin/releases) section of this repository
   2. Install the plugin with `cf install-plugin <path_to_binary>`. Use -f flag to uninstall existing plugin if any and install the new one. 
   
**Option 2:**

If you are using MacOS, you could run  

    cf install-plugin -f https://github.com/faisaltheparttimecoder/cf-cli-delete-wrapper-plugin/releases/download/0.1.1/cf-delete-wrapper_v0.1.1.osx
    
# Usage

1. **Command:** delete-multi-apps, **Alias:** dma 
    ```
    $ cf delete-multi-apps --help
      NAME:
         delete-multi-apps - Delete multiple apps via a single command
      
      ALIAS:
         dma
      
      USAGE:
         cf delete-multi-apps -a <APP1>,<APP2>,....,<APPn>
      
      OPTIONS:
         -force       -f, no need to prompt for confirmation
         -app         -a, list of apps to be deleted
    ```

2. **Command:** delete-app-using-manifest, **Alias:** daum
    ```
    $ cf delete-app-using-manifest --help
    NAME:
       delete-app-using-manifest - Detect the apps name from manifest and delete it
    
    ALIAS:
       daum
    
    USAGE:
       cf delete-app-using-manifest
    
    OPTIONS:
       -force       -f, no need to prompt for confirmation
    ```
    
# Example

+ To delete multiple at once
    ```
    $ cf dma -a test1,test2,test3
    Are you sure you want to delete these apps (test1,test2,test3), do you wish to continue (Yy/Nn)?: y
    
    Successfully deleted the app "test1"
    
    Successfully deleted the app "test2"
    
    Successfully deleted the app "test3"
    ```

+ To delete an app that is on the manifest
    ```
    $ cat manifest.yml 
      ---
      applications:
      - name: customer-test
    
    $ cf daum
    Are you sure you want to delete these apps (customer-test), do you wish to continue (Yy/Nn)?: y
    
    Successfully deleted the app "customer-test"
    ```

# Build

```
go get github.com/faisaltheparttimecoder/cf-cli-delete-wrapper-plugin/
-- Modify the code
run "/bin/sh run.sh" to build package
```