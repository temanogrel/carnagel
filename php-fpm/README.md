# PHP-FPM role

## Core packages

The following packages are considered core and should not be overwritten since they will always reflect
the base required to run zend framework/expressive
```
php_packages:
  - php{{ php_version }}-cli
  - php{{ php_version }}-common
  - php{{ php_version }}-dev
  - php{{ php_version }}-curl
  - php{{ php_version }}-fpm
  - php{{ php_version }}-intl
  - php{{ php_version }}-json
  - php{{ php_version }}-xml
  - php{{ php_version }}-mbstring
```

## Extended packages

To add custom packages to an installation just use the `php_extended_packages` instead

```
php_extended_packages:
  - php{{ php_version }}-bcmath
```

## Select database driver(s)
```
php_mysql_enabled: false
php_sqlite_enabled: false
php_postgresql_enabled: false
```

## Xdebug

The `xdebug_remote_host_ip` should be configured to the the first ip of the used private subnet in vagrant, so if your
vagrant uses the private network `192.168.10.x` then you should set it to `192.168.10.1`.

```
xdebug_enabled: false
xdebug_branch: master
xdebug_remote_host_ip: false
```

## Redis
```
redis_enabled: false
redis_branch: php7
```

## Newrelic
```
newrelic_php_enabled: 1
newrelic_project_name: 'hello world'
newrelic_license_key: 'the key'
```
