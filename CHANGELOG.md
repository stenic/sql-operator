# Changelog

All notable changes to this project will be documented in this file. See
[Conventional Commits](https://conventionalcommits.org) for commit guidelines.

## [1.12.1](https://github.com/stenic/sql-operator/compare/v1.12.0...v1.12.1) (2021-11-15)


### Bug Fixes

* Add ServiceMonitor ([05e9af0](https://github.com/stenic/sql-operator/commit/05e9af0c16615579e085e48f3be2742039e06b1a))

# [1.12.0](https://github.com/stenic/sql-operator/compare/v1.11.0...v1.12.0) (2021-11-15)


### Features

* Add operation metrics ([7552bb4](https://github.com/stenic/sql-operator/commit/7552bb4365111d2ccd8e11606951ed4f54bdc92a))

# [1.11.0](https://github.com/stenic/sql-operator/compare/v1.10.1...v1.11.0) (2021-11-14)


### Features

* Change SqlHost to DSN ([a7c3872](https://github.com/stenic/sql-operator/commit/a7c3872964d4689ad826af7416ff90403241dcd9)), closes [#37](https://github.com/stenic/sql-operator/issues/37)
* Change SqlHost to DSN ([7613db7](https://github.com/stenic/sql-operator/commit/7613db7b9395d8c14b388ed257482aaaf3f5c243)), closes [#37](https://github.com/stenic/sql-operator/issues/37)

## [1.10.1](https://github.com/stenic/sql-operator/compare/v1.10.0...v1.10.1) (2021-11-14)


### Bug Fixes

* Enable leader-elect when running ha ([43d306c](https://github.com/stenic/sql-operator/commit/43d306c1e81fc8aaa1d6768f86477be1161d6d81))

# [1.10.0](https://github.com/stenic/sql-operator/compare/v1.9.1...v1.10.0) (2021-11-14)


### Bug Fixes

* **deps:** update module github.com/onsi/gomega to v1.17.0 ([61f774e](https://github.com/stenic/sql-operator/commit/61f774ef9509974d72fa87bed887db7cf486255e))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.10.3 ([b1de5fb](https://github.com/stenic/sql-operator/commit/b1de5fbb429b8f0d3f4d8dc13f349953981df47d))
* Disable webhook by default ([f8f2e85](https://github.com/stenic/sql-operator/commit/f8f2e851a0d78fdd571f95682b3a52f3e33a4370))
* Publish driver errors to events ([d4cfb17](https://github.com/stenic/sql-operator/commit/d4cfb1781c632af8a17fc5899320d3307e58c05d)), closes [#38](https://github.com/stenic/sql-operator/issues/38)


### Features

* Don't delete object befor it's children are deleted ([7812dcb](https://github.com/stenic/sql-operator/commit/7812dcbc83863f1e8cea5d13d43fe2949a46e29b)), closes [#33](https://github.com/stenic/sql-operator/issues/33)

## [1.9.1](https://github.com/stenic/sql-operator/compare/v1.9.0...v1.9.1) (2021-11-05)


### Bug Fixes

* Cleanup debug statements ([09b568e](https://github.com/stenic/sql-operator/commit/09b568e789554df7ee3ea66a7191be5f3b37e4a8))

# [1.9.0](https://github.com/stenic/sql-operator/compare/v1.8.2...v1.9.0) (2021-11-05)


### Features

* Implement delete grants ([c1c7694](https://github.com/stenic/sql-operator/commit/c1c76941aca6fa90e7011e749613ee58b5aa501c))

## [1.8.2](https://github.com/stenic/sql-operator/compare/v1.8.1...v1.8.2) (2021-11-05)


### Bug Fixes

* Handle missing secret ([6ecd51d](https://github.com/stenic/sql-operator/commit/6ecd51d856ddddf06d6f9d2ba2889a37b5497eed))

## [1.8.1](https://github.com/stenic/sql-operator/compare/v1.8.0...v1.8.1) (2021-11-05)


### Bug Fixes

* Handle missing secret ([09fe2f8](https://github.com/stenic/sql-operator/commit/09fe2f8d9b2abfc046eeecae2e728f336475b404))

# [1.8.0](https://github.com/stenic/sql-operator/compare/v1.7.0...v1.8.0) (2021-11-05)


### Features

* Disable webhooks ([9c56acc](https://github.com/stenic/sql-operator/commit/9c56acc4bbe9e75bc340c2f834465dd8d44c39e4))

# [1.7.0](https://github.com/stenic/sql-operator/compare/v1.6.3...v1.7.0) (2021-11-05)


### Features

* Allow to set log-level ([0831998](https://github.com/stenic/sql-operator/commit/083199843d274ff236a8555f4c4b2740931c7913))
* Allow to set log-level ([5833b6d](https://github.com/stenic/sql-operator/commit/5833b6d4fff5ce0ab414e81fee424c9949a29c19))

## [1.6.3](https://github.com/stenic/sql-operator/compare/v1.6.2...v1.6.3) (2021-11-05)


### Bug Fixes

* Don't limit build target ([edcfba5](https://github.com/stenic/sql-operator/commit/edcfba50725159ee681823744fb78e5215d0b4d7))

## [1.6.2](https://github.com/stenic/sql-operator/compare/v1.6.1...v1.6.2) (2021-11-05)


### Bug Fixes

* Remove debug statement ([e9c5ab9](https://github.com/stenic/sql-operator/commit/e9c5ab94ea409fb6c0b506ed16d6ea14cf856774))

## [1.6.1](https://github.com/stenic/sql-operator/compare/v1.6.0...v1.6.1) (2021-11-05)


### Bug Fixes

* Handle grant lookup better ([70b9b25](https://github.com/stenic/sql-operator/commit/70b9b25ce83f3d7602ab6d4ed4395cc497f1e3f6))

# [1.6.0](https://github.com/stenic/sql-operator/compare/v1.5.19...v1.6.0) (2021-11-04)


### Features

* Add events to objects ([5c970f9](https://github.com/stenic/sql-operator/commit/5c970f9a047f5fadda693150c057377bfc9f1bd2))

## [1.5.19](https://github.com/stenic/sql-operator/compare/v1.5.18...v1.5.19) (2021-11-02)


### Bug Fixes

* Remove host column from SqlGrant list ([deac535](https://github.com/stenic/sql-operator/commit/deac535709fcd5162a4db0df4ed5b006382f8184))

## [1.5.18](https://github.com/stenic/sql-operator/compare/v1.5.17...v1.5.18) (2021-11-01)


### Bug Fixes

* Correct webhook naming ([dc35b32](https://github.com/stenic/sql-operator/commit/dc35b32008c8ef69ca9f6e14632a0f0a224a2dee))

## [1.5.17](https://github.com/stenic/sql-operator/compare/v1.5.16...v1.5.17) (2021-11-01)


### Bug Fixes

* Less history should be fine ([30ea82a](https://github.com/stenic/sql-operator/commit/30ea82a19efaaf57031271b0c3283b31facd8449))

## [1.5.16](https://github.com/stenic/sql-operator/compare/v1.5.15...v1.5.16) (2021-11-01)


### Bug Fixes

* Correct prometheus scrape port ([989b1f8](https://github.com/stenic/sql-operator/commit/989b1f87ac58ef4375427a5a63f290971962a969))

## [1.5.15](https://github.com/stenic/sql-operator/compare/v1.5.14...v1.5.15) (2021-10-31)


### Bug Fixes

* **deps:** Patch github.com/onsi/ginkgo ([3b4d5a3](https://github.com/stenic/sql-operator/commit/3b4d5a33aa77359285cde1554359529f2888043a))
* **deps:** Patch github.com/onsi/gomega ([d076abb](https://github.com/stenic/sql-operator/commit/d076abbee13a11260edd882e6b33a32a7bb25048))
* **deps:** Patch k8s.io ([29df11f](https://github.com/stenic/sql-operator/commit/29df11f8748963c7076b48bc55c27de40a83ac22))

## [1.5.14](https://github.com/stenic/sql-operator/compare/v1.5.13...v1.5.14) (2021-10-31)


### Bug Fixes

* Add status badge ([e83d684](https://github.com/stenic/sql-operator/commit/e83d684c4743c0c1922b1856a68cbfd953105804))

## [1.5.13](https://github.com/stenic/sql-operator/compare/v1.5.12...v1.5.13) (2021-10-31)


### Reverts

* Revert "chore(deps): Update vulnerable dependecies" ([b229129](https://github.com/stenic/sql-operator/commit/b229129711515e8ab446a342de58e97751a1e6a9))
* Revert "chore(deps): update golang docker tag to v1.17" ([8928d75](https://github.com/stenic/sql-operator/commit/8928d759fabf49083940c595e5ff2ffcd54a1b2f))

## [1.5.12](https://github.com/stenic/sql-operator/compare/v1.5.11...v1.5.12) (2021-10-30)


### Bug Fixes

* Limit db connections ([8a870a2](https://github.com/stenic/sql-operator/commit/8a870a28b775f5be768d7b5172a5c943e3d148e4))

## [1.5.11](https://github.com/stenic/sql-operator/compare/v1.5.10...v1.5.11) (2021-10-30)


### Bug Fixes

* **mysql:** Use exec instead of query if possible ([ddfe3f1](https://github.com/stenic/sql-operator/commit/ddfe3f14e0ae43cb52f7b1fa8303840c26f30402))

## [1.5.10](https://github.com/stenic/sql-operator/compare/v1.5.9...v1.5.10) (2021-10-30)


### Bug Fixes

* **chart:** No leader-elect ([e78b137](https://github.com/stenic/sql-operator/commit/e78b1375847bc799f30951bde6b87e4fd24e254b))

## [1.5.9](https://github.com/stenic/sql-operator/compare/v1.5.8...v1.5.9) (2021-10-30)


### Bug Fixes

* **chart:** Handle appVersion from chart ([3246146](https://github.com/stenic/sql-operator/commit/3246146ac60599891130726892e325503cc11085))

## [1.5.8](https://github.com/stenic/sql-operator/compare/v1.5.7...v1.5.8) (2021-10-30)


### Bug Fixes

* **chart:** This is now set on release ([ff56cb0](https://github.com/stenic/sql-operator/commit/ff56cb03b69b3bfba9e2175a6a8a5e4cdf1503c8))

## [1.5.7](https://github.com/stenic/sql-operator/compare/v1.5.6...v1.5.7) (2021-10-30)


### Bug Fixes

* **ci:** Correct chart location ([d18accf](https://github.com/stenic/sql-operator/commit/d18accf88fba0f0448d41373c2cce8723a20cdda))
* **ci:** Remove config ([749a560](https://github.com/stenic/sql-operator/commit/749a560acf2d6145aa5ed84f83d62c2af7b5b7d5))
* Minor code improvements ([2fc9efc](https://github.com/stenic/sql-operator/commit/2fc9efcbd81b3882c31570e51d90d4f401419eb3))

## [1.5.6](https://github.com/stenic/sql-operator/compare/v1.5.5...v1.5.6) (2021-10-30)


### Bug Fixes

* Minor code improvements ([4c6d05f](https://github.com/stenic/sql-operator/commit/4c6d05f6b5d0ef6814b420367034f47493ccade0))

## [1.5.5](https://github.com/stenic/sql-operator/compare/v1.5.4...v1.5.5) (2021-10-29)


### Bug Fixes

* Minor code improvements ([88f9a88](https://github.com/stenic/sql-operator/commit/88f9a884058f04612e728e6f08c2d145a4476df0))

## [1.5.4](https://github.com/stenic/sql-operator/compare/v1.5.3...v1.5.4) (2021-10-29)


### Bug Fixes

* Minor code improvements ([2b6c99f](https://github.com/stenic/sql-operator/commit/2b6c99f7274a22e2a6627c7a61290ffa5b31cf2a))

## [1.5.3](https://github.com/stenic/sql-operator/compare/v1.5.2...v1.5.3) (2021-10-29)


### Bug Fixes

* Minor code improvements ([7ad6360](https://github.com/stenic/sql-operator/commit/7ad6360008d28c72dafa0ef7c69c64baf76ef2ec))

## [1.5.2](https://github.com/stenic/sql-operator/compare/v1.5.1...v1.5.2) (2021-10-29)


### Bug Fixes

* Minor code improvements ([6fc862c](https://github.com/stenic/sql-operator/commit/6fc862c0b0918884beeccba708adef2a3fc89418))

## [1.5.1](https://github.com/stenic/sql-operator/compare/v1.5.0...v1.5.1) (2021-10-29)


### Bug Fixes

* Cleanup ([75169ce](https://github.com/stenic/sql-operator/commit/75169ce43a0dea31284219230f6c483e9aa0ee15))

# [1.5.0](https://github.com/stenic/sql-operator/compare/v1.4.0...v1.5.0) (2021-10-28)


### Features

* No next ([1a1fa31](https://github.com/stenic/sql-operator/commit/1a1fa31d8388069f80e3684ddbce4756157b615d))

# [1.4.0](https://github.com/stenic/sql-operator/compare/v1.3.1...v1.4.0) (2021-10-28)


### Features

* Add next for next release ([f21709f](https://github.com/stenic/sql-operator/commit/f21709fce7b608d643901c2616b4afbb33b95996))
* No need for sqlHost on SqlGrant ([3e8c51b](https://github.com/stenic/sql-operator/commit/3e8c51bd5ab8104555bda132a7e97cb0bfe9be12))

## [1.3.1](https://github.com/stenic/sql-operator/compare/v1.3.0...v1.3.1) (2021-10-25)


### Bug Fixes

* **chart:** Fix missing role ([2a27690](https://github.com/stenic/sql-operator/commit/2a276905a661206571bf539779b7edbedde2c8a6))
* **chart:** Fix missing role ([6dd7de5](https://github.com/stenic/sql-operator/commit/6dd7de5d54be4da56779e6a18983e391a738e212))

# [1.3.0](https://github.com/stenic/sql-operator/compare/v1.2.4...v1.3.0) (2021-10-25)


### Bug Fixes

* Restore boilerplate location ([acd9337](https://github.com/stenic/sql-operator/commit/acd93372d059c0acfb12606b96553ac753568c4e))


### Features

* Make fields immutable [#7](https://github.com/stenic/sql-operator/issues/7) ([8c19f80](https://github.com/stenic/sql-operator/commit/8c19f802f125f5e50500a20d528929b6575cd982))
* Validate datbasename [#6](https://github.com/stenic/sql-operator/issues/6) ([97431a0](https://github.com/stenic/sql-operator/commit/97431a015f1273cdcf7da0d1c41f012818ea26d4))

## [1.2.4](https://github.com/stenic/sql-operator/compare/v1.2.3...v1.2.4) (2021-10-23)


### Bug Fixes

* Render README.md ([f78f6d4](https://github.com/stenic/sql-operator/commit/f78f6d4b00914389e0bf11a6dc7f298a9fe4baae))

## [1.2.3](https://github.com/stenic/sql-operator/compare/v1.2.2...v1.2.3) (2021-10-23)


### Bug Fixes

* Render README.md ([bfc75f0](https://github.com/stenic/sql-operator/commit/bfc75f033d6b5dcf47b6d46a84d7ca31cd7d1088))
* Render README.md ([55c5a20](https://github.com/stenic/sql-operator/commit/55c5a2044dac5d82bc02c2d5f95dcb1fe3e27e92))

## [1.2.2](https://github.com/stenic/sql-operator/compare/v1.2.1...v1.2.2) (2021-10-23)


### Bug Fixes

* Fix copy-past error ([46fdcdd](https://github.com/stenic/sql-operator/commit/46fdcdd616dc1a38f464b6979bbc4514306be976))

## [1.2.1](https://github.com/stenic/sql-operator/compare/v1.2.0...v1.2.1) (2021-10-23)


### Bug Fixes

* Add chart README.md ([3bb8129](https://github.com/stenic/sql-operator/commit/3bb812992677b0a74dfa688583426d64f77d14b6))

# [1.2.0](https://github.com/stenic/sql-operator/compare/v1.1.0...v1.2.0) (2021-10-23)


### Bug Fixes

* Cleanup ([3b3cabd](https://github.com/stenic/sql-operator/commit/3b3cabd7687b9a1c6a2cabc171ce0293a40ec326))


### Features

* Release charts ([55a36b6](https://github.com/stenic/sql-operator/commit/55a36b6622e3358a2e3360da5fdd0057585ed2fd))
* Release charts ([6054a32](https://github.com/stenic/sql-operator/commit/6054a32a076ce590be0edc15896967ef8f9eb76d))
* Update chart maintainers ([84cea14](https://github.com/stenic/sql-operator/commit/84cea14d741aaad3f596edb7129804b9f20b60d8))

# [1.1.0](https://github.com/stenic/sql-operator/compare/v1.0.0...v1.1.0) (2021-10-23)


### Features

* Add pipeline ([#5](https://github.com/stenic/sql-operator/issues/5)) ([029729c](https://github.com/stenic/sql-operator/commit/029729cf639a88fc61fbf478390413b09938c34a))
* More builds ([031f8c5](https://github.com/stenic/sql-operator/commit/031f8c594cd28f8cc09dd7835fa4167e9b5bc289))

# 1.0.0 (2021-10-23)


### Features

* Add database ([5fdf098](https://github.com/stenic/sql-operator/commit/5fdf098a5c0acf1ef3d7f21cacaf36fa7b14bba3))
* First draft ([cd63976](https://github.com/stenic/sql-operator/commit/cd63976b46300416de891ad7f412f7674cee7344))
* Implement create user ([144d23a](https://github.com/stenic/sql-operator/commit/144d23ae42b0c5568a356a57070d3991c1a959a2))
* Implement grants ([a45ce70](https://github.com/stenic/sql-operator/commit/a45ce708c1b4f9f8d5c2d59946659e4428bf7104))
* Init controlloop ([3f6614b](https://github.com/stenic/sql-operator/commit/3f6614b5ac0b7da268e0e77139d680f291f7c7d5))
* Only apply grant if needed ([81429a7](https://github.com/stenic/sql-operator/commit/81429a79a9751338fa91fb51a73099e978f68456))
