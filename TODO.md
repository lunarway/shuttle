# Validating the use-case of shuttle

## Services we should test

* [apollo] [nodejs] travelcard -
  * docker build ✅
  * full jenkins run 🚀
    * point to db?
    * point to secrets?
    * shared env specs between environments?
  * manage secrets?
* [wannabe-nasa] [nodejs] support2 - https://bitbucket.org/LunarWay/lunar-way-support-service
* [nasa] [go] subscription - https://bitbucket.org/LunarWay/lunar-way-subscription-service
* [nasa] [go] synchronize - https://bitbucket.org/LunarWay/lunar-way-synchronize-service
* [nasa] [go] localization - https://bitbucket.org/LunarWay/lunar-way-localization-service
* [nasa] [go] kpi - https://bitbucket.org/LunarWay/lunar-way-kpi-service
* [nasa] [go] authentication - https://bitbucket.org/LunarWay/lunar-way-authentication-service
* [nasa] [go] bec-exporter - https://bitbucket.org/LunarWay/lunar-way-bec-service
* [nasa] [go] user - https://bitbucket.org/LunarWay/lunar-way-user-service
* [nasa] [go] appsync - https://bitbucket.org/LunarWay/lunar-way-appsync-service
* [nasa] [go] upodi & upodi-monitoring - https://bitbucket.org/LunarWay/lunar-way-upodi-service
* [nasa] [go] intercom-sync - https://bitbucket.org/LunarWay/lunar-way-intercom-service
* [nasa] [go] product - https://bitbucket.org/LunarWay/lunar-way-product-service


# Features

* tmp directory
* overload templates in project
* add scripts in project & overload scripts
* template functions:
  * Constant string formatting + other likewise
* lock to git version of plan
  * update on run
  * CI update plan git coupling
* Temporarily Point to local plan instead of git for local development

## Less needed features, but still awesome

* schema for vars
* scripts arguments with flags as alternative to positional