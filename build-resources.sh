#!/bin/bash

files="files"
vendor="files/js/vendor"
build="files/build"

compile=$(echo "${vendor}/jquery-2.2.4.js" \
"${vendor}/jquery.highlight.js" \
"${vendor}/nprogress.js" \
"${vendor}/jquery.datetimepicker.js" \
"${vendor}/doT.js" \
"${vendor}/jquery.doT.plugin.js" \
"${vendor}/jquery.easing.1.3.min.js" \
"${vendor}/jquery.switch.js" \
"${vendor}/jquery.tabs.js" \
"${vendor}/jquery.message.js" \
"${vendor}/jquery.tools.js" \
"${vendor}/redactor.js" \
"${vendor}/jquery.shwark.js" \
"${vendor}/jquery.cookie.js" \
"${vendor}/groupper.js" \
"${vendor}/jquery.wbox.js" \
"${vendor}/jquery.ajaxHelper.js" \
"${vendor}/bootstrap.min.js" \
"${vendor}/select2.full.js" \
"${vendor}/jquery.perfect-scrollbar.min.js" \
"${vendor}/moment-with-locales.min.js" \
"${vendor}/list.js")

compileLogin=$(echo "${vendor}/jquery-2.2.4.js" \
"${vendor}/jquery.easing.1.3.min.js" \
"${vendor}/jquery.message.js" \
"${vendor}/nprogress.js" \
"${vendor}/jquery.ajaxHelper.js")

uglifyjs --source-map filename="${build}/main.js.map" --source-map-root "/" --source-map-url "/${build}/main.js.map" -o "${build}/main.js" $compile "files/js/common.js"

uglifyjs --source-map filename="${build}/login.js.map" --source-map-root "/" --source-map-url "/${build}/login.js.map" -o "${build}/login.js" $compileLogin

cd "files/less"
lessc --clean-css="--s1 --advanced --compatibility=ie8" style.less > ../build/style.css
lessc --clean-css="--s1 --advanced --compatibility=ie8" login.less > ../build/login.css
cd ../..

gzip "${build}/main.js" -c > "${build}/main.js.gz"
gzip "${build}/login.js" -c > "${build}/login.js.gz"
gzip "${build}/style.css" -c > "${build}/style.css.gz"
gzip "${build}/login.css" -c > "${build}/login.css.gz"
