echo "git commit -a -m '第v1.4.$1版本'"
git commit -a -m "第v1.4.$1版本"
echo "git checkout -b release/v1.4.$1"
git checkout -b release/v1.4.$1
echo "git push -u origin release/v1.4.$1"
git push -u origin release/v1.4.$1
echo "git tag v1.4.$1"
git tag v1.4.$1
echo "git push --tags"
git push --tags