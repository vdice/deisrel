# Change Log

## [Unreleased](https://github.com/deis/deisrel/tree/HEAD)

**Implemented enhancements:**

- create command for generating changelogs for an individual repository [\#45](https://github.com/deis/deisrel/issues/45)
- generate-changelog: Should output what repo the message came from  [\#43](https://github.com/deis/deisrel/issues/43)
- publish deisrel binary to bintray [\#18](https://github.com/deis/deisrel/issues/18)
- Pull down workflow-dev\(-e2e\) charts to gen release versions [\#13](https://github.com/deis/deisrel/issues/13)
- Create command for generating changelogs [\#5](https://github.com/deis/deisrel/issues/5)
- Create command for tagging repos [\#3](https://github.com/deis/deisrel/issues/3)
- Add command for aggregating all changelogs [\#2](https://github.com/deis/deisrel/issues/2)
- Automate release chart\(s\) creation [\#1](https://github.com/deis/deisrel/issues/1)
- feat\(getShas\): add ability to pass in ref to git shas cmd [\#56](https://github.com/deis/deisrel/pull/56) ([vdice](https://github.com/vdice))
- feat\(.github\): add pr template with release checklist reminder [\#50](https://github.com/deis/deisrel/pull/50) ([vdice](https://github.com/vdice))
- fix\(main.go\): change GITHUB\_TOKEN to GITHUB\_ACCESS\_TOKEN [\#48](https://github.com/deis/deisrel/pull/48) ([arschles](https://github.com/arschles))
- feat\(main.go,actions\): add command to generate a changelog entry for an individual repo [\#47](https://github.com/deis/deisrel/pull/47) ([arschles](https://github.com/arschles))
- fix\(actions/generate\_changelog.go\): output the repo for each commit message [\#46](https://github.com/deis/deisrel/pull/46) ([arschles](https://github.com/arschles))
- feat\(ci\): properly update app version, update build [\#30](https://github.com/deis/deisrel/pull/30) ([vdice](https://github.com/vdice))
- feat\(generate\_params.go\): add versionsApiURL to prod for release chart [\#29](https://github.com/deis/deisrel/pull/29) ([vdice](https://github.com/vdice))
- feat\(actions\): add deisrel git tag command [\#26](https://github.com/deis/deisrel/pull/26) ([bacongobbler](https://github.com/bacongobbler))
- feat\(helm\_generate\): enable staging of generated params [\#24](https://github.com/deis/deisrel/pull/24) ([vdice](https://github.com/vdice))
- feat\(actions\): add doc\( tag for generating changelogs [\#21](https://github.com/deis/deisrel/pull/21) ([bacongobbler](https://github.com/bacongobbler))
- feat\(ci/release\): publish binary to bintray [\#20](https://github.com/deis/deisrel/pull/20) ([vdice](https://github.com/vdice))
- feat\(ci\): add Makefile and .travis.yml [\#19](https://github.com/deis/deisrel/pull/19) ([vdice](https://github.com/vdice))
- docs\(badge\): added code-beat badge [\#17](https://github.com/deis/deisrel/pull/17) ([chaitanyaenr](https://github.com/chaitanyaenr))
- feat\(actions\): add generate-changelog action [\#16](https://github.com/deis/deisrel/pull/16) ([bacongobbler](https://github.com/bacongobbler))

**Fixed bugs:**

- GITHUB\_TOKEN should be GITHUB\_ACCESS\_TOKEN [\#41](https://github.com/deis/deisrel/issues/41)
- fix\(actions/release\_walker.go\): fix issue where /dev/null was changed [\#67](https://github.com/deis/deisrel/pull/67) ([vdice](https://github.com/vdice))

**Closed issues:**

- implement proposal \#52 [\#54](https://github.com/deis/deisrel/issues/54)
- \[proposal\] refactor helm-params,stage [\#52](https://github.com/deis/deisrel/issues/52)
- Refactor global changelog generator to use changelog package functions [\#51](https://github.com/deis/deisrel/issues/51)
- create documentation on how to use the tool [\#39](https://github.com/deis/deisrel/issues/39)
- need to update cli [\#33](https://github.com/deis/deisrel/issues/33)
- document deisrel [\#15](https://github.com/deis/deisrel/issues/15)

**Merged pull requests:**

- fix\(get\_shas.go\): ref wasn't getting to get shas method [\#65](https://github.com/deis/deisrel/pull/65) ([vdice](https://github.com/vdice))
- fix\(actions/common.go\): remove 'workflow' from repo/component list... [\#64](https://github.com/deis/deisrel/pull/64) ([vdice](https://github.com/vdice))
- chore\(generate\_params.go\): set default doctor api url to prod [\#63](https://github.com/deis/deisrel/pull/63) ([vdice](https://github.com/vdice))
- docs\(README\): document deisrel [\#61](https://github.com/deis/deisrel/pull/61) ([bacongobbler](https://github.com/bacongobbler))
- ref\(helm\_stage/generate\): collapse multiple helm chart commands into one [\#55](https://github.com/deis/deisrel/pull/55) ([vdice](https://github.com/vdice))
- ref\(actions,changelog\): use changelog package in global changelog command [\#53](https://github.com/deis/deisrel/pull/53) ([arschles](https://github.com/arschles))
- feat\(Makefile\): append VERSION to binary name for uniqueness [\#44](https://github.com/deis/deisrel/pull/44) ([vdice](https://github.com/vdice))
- feat\(helm\_stage\_e2e.go\): add 'tpl/workflow-e2e-pod.yaml' to list [\#36](https://github.com/deis/deisrel/pull/36) ([vdice](https://github.com/vdice))
- fix\(helm\_generate e2e\): fix helm generate e2e files [\#35](https://github.com/deis/deisrel/pull/35) ([vdice](https://github.com/vdice))
- fix\(\*\): update action func to adhere to codegangsta 1.15.0 spec [\#34](https://github.com/deis/deisrel/pull/34) ([sgoings](https://github.com/sgoings))
- fix\(helm\_generate\): add staging code lost in auto-merge [\#32](https://github.com/deis/deisrel/pull/32) ([vdice](https://github.com/vdice))
- chore\(RELEASE vars\): update DEIS -\> WORKFLOW [\#31](https://github.com/deis/deisrel/pull/31) ([vdice](https://github.com/vdice))
- chore\(params\): missed adding stdout metrics to the params file [\#25](https://github.com/deis/deisrel/pull/25) ([jchauncey](https://github.com/jchauncey))
- feat\(common\): Add stdout-metrics [\#23](https://github.com/deis/deisrel/pull/23) ([jchauncey](https://github.com/jchauncey))
- fix\(bintray-template.json\): subject change to 'deis' [\#22](https://github.com/deis/deisrel/pull/22) ([vdice](https://github.com/vdice))
- feat\(helm-chart e2e/workflow\): grab latest version\(s\) of workflow dev… [\#14](https://github.com/deis/deisrel/pull/14) ([vdice](https://github.com/vdice))
- feat\(helm\_generate\_workflow\): can now generate workflow params [\#10](https://github.com/deis/deisrel/pull/10) ([vdice](https://github.com/vdice))
- ref\(misc\): rename arschles -\> deis, generage -\> generate [\#7](https://github.com/deis/deisrel/pull/7) ([vdice](https://github.com/vdice))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*