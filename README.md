# Usage

## Execution

run:
```
make run
```

Also parameters host and port for the queue can also be changed (hereinafter **run-options**):

```
*only host*  
make run host=somehost
*only port*  
make run port=123
*both*  
make run host=somehost port=123
```

build:
```
make build
```

build and run:
```
make build-run <*run-options...*>
```

run already builded earlier:
```
make run-builded <*run-options...*>
```

## Commands

Before each command you can see list of current options:  
  
![Alt text](/.documentation/images/options.png "List of options")  
  

1. **Show list of active tubes in queue**
```
list
```

2. **Change options**

- Changing single option:
```
set tube=$name
set delay=$value
set ttr=$value
set priority=$value
```

- Changing multiple options via a single set:
```
set tube=$name delay=$value ttr=$value priority=$value
set tube=$name priority=$value
*etc...*
```

3. **Put data**

- From console
```
put -d data...
```

- From file
```
put -f path/to/file
```