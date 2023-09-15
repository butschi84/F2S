#!/bin/bash
echo $1

git pull
git merge origin/helm-release-$1
git push origin --delete f2s-$1
git push origin --delete helm-release-$1
git tag $1
git tag -d f2s-$1
git push
git push origin $1