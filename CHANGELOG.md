# Changelog

## [1.2.2](https://github.com/meysam81/liveness-check/compare/v1.2.1...v1.2.2) (2025-06-15)


### Bug Fixes

* list the pods again if the pod ip is not yet assigned ([444d690](https://github.com/meysam81/liveness-check/commit/444d6900ed0f687e27ef6056f8c1eebdd3fb815f))

## [1.2.1](https://github.com/meysam81/liveness-check/compare/v1.2.0...v1.2.1) (2025-06-15)


### Bug Fixes

* **docs:** add more badges ([6e98d97](https://github.com/meysam81/liveness-check/commit/6e98d97593273468832104b4a107c151cd761b5d))
* **docs:** remove extra backticks ([dd60b54](https://github.com/meysam81/liveness-check/commit/dd60b54699df2bd2a92152884091cbea283d1cbd))

## [1.2.0](https://github.com/meysam81/liveness-check/compare/v1.1.2...v1.2.0) (2025-06-15)


### Features

* add k8s pod checker ([2d3ae8f](https://github.com/meysam81/liveness-check/commit/2d3ae8f48818f3a38e1ddd21157e6d8384f83954))
* **docs:** beautify and provide clear examples ([2abf984](https://github.com/meysam81/liveness-check/commit/2abf9842aaf15dde4de60b91376d907b4a310a49))


### Bug Fixes

* **build:** increase build by caching docker layers ([448c89d](https://github.com/meysam81/liveness-check/commit/448c89d445345bedc27edb959dbc538832339f99))
* **CI:** do not check compare URLs ([2188f41](https://github.com/meysam81/liveness-check/commit/2188f41fb1f91ee807cfcbb367bb796216bccecd))
* **docs:** add line separator to headings ([8501b03](https://github.com/meysam81/liveness-check/commit/8501b038432e1ff1b54bb7a092bc5a26c4209362))
* **docs:** modify badges ([1ce0c05](https://github.com/meysam81/liveness-check/commit/1ce0c05dc60e0b6387c954d122f2db9584b0f0f6))
* **docs:** provide a complete k8s example ([63069fe](https://github.com/meysam81/liveness-check/commit/63069fe734e4c9a457817213635f75d523987282))
* quickly exit when the max retries is reached ([49f462c](https://github.com/meysam81/liveness-check/commit/49f462c20a36b2387b097c3774ed8cd0285b1ca4))
* return error on listing pods ([75c4551](https://github.com/meysam81/liveness-check/commit/75c4551b29f2d875b947c21cf93dcffef45c9616))
* use ticker instead of time and lower log level ([8ee8361](https://github.com/meysam81/liveness-check/commit/8ee83612773de2c4c1297df83d09788d8a3bda05))

## [1.1.2](https://github.com/meysam81/liveness-check/compare/v1.1.1...v1.1.2) (2025-06-14)


### Bug Fixes

* **build:** append changelog ([cbbf307](https://github.com/meysam81/liveness-check/commit/cbbf307e1bc829daa8a2eb619713f5339e1e20b1))
* **docs:** add docker pull count ([d068bfb](https://github.com/meysam81/liveness-check/commit/d068bfb32da713a68db5700afe94059c5f477aa1))
* **docs:** add k8s example ([526df15](https://github.com/meysam81/liveness-check/commit/526df153b50a9d06ae50a02eb5a66ced7f0d8798))
* **docs:** run k8s job as init container ([e857219](https://github.com/meysam81/liveness-check/commit/e857219ffefb9b8b09a0ae0e134edf3d15bd5782))
* drop log level from cli args ([30a2992](https://github.com/meysam81/liveness-check/commit/30a2992c977e727fd1ef876f941f43b04b0c00a7))
* initialize log level manually from envs ([98d862d](https://github.com/meysam81/liveness-check/commit/98d862d5fa6340bc49aab9121bddd432078a5447))

## [1.1.1](https://github.com/meysam81/liveness-check/compare/v1.1.0...v1.1.1) (2025-06-14)


### Bug Fixes

* **build:** uncomment the cosign ([a32cfe9](https://github.com/meysam81/liveness-check/commit/a32cfe9e2a30f6a0a7bba57189cf44f7b4a1660c))
* perform major refactor in favor of modularity and maintainability ([9bf8017](https://github.com/meysam81/liveness-check/commit/9bf8017b9e209b0ca0b50fb91a54e22d67f8d09c))

## [1.1.0](https://github.com/meysam81/liveness-check/compare/v1.0.0...v1.1.0) (2025-06-14)


### Features

* add goreleaser ([b921b70](https://github.com/meysam81/liveness-check/commit/b921b70dfbca6baac366f2c3aac832ab1b58f0bb))
* handle shutdown signal ([c1d79b5](https://github.com/meysam81/liveness-check/commit/c1d79b541451c85ba9cf0c37106f9781be895c3e))


### Bug Fixes

* **CI:** do not fail on empty links ([e26d7af](https://github.com/meysam81/liveness-check/commit/e26d7afbb9003d3569a45b7ae0029669e546a435))
* **dev:** enable golangci pre-commit hook ([b9cbde2](https://github.com/meysam81/liveness-check/commit/b9cbde2d6640efd98edbc243e3551a1a2d78329c))
* handle version within urfave ([5a6527e](https://github.com/meysam81/liveness-check/commit/5a6527e521f3e39bb0502a36ab8f5964f102d550))
* ignore secrets ([cd394b8](https://github.com/meysam81/liveness-check/commit/cd394b8e87249d2898d76e906d3c96009e0155b6))

## 1.0.0 (2025-06-14)


### Features

* add the go application ([83cf257](https://github.com/meysam81/liveness-check/commit/83cf257a22d7f2a65be2f429a8b877fe7b77df19))
* add version subcommand and move check to its own ([fdc35ce](https://github.com/meysam81/liveness-check/commit/fdc35ced3ca30ae02e54960d9a04d35fb6f8c24c))


### Bug Fixes

* do not stop unless retries is set ([d8f0bdf](https://github.com/meysam81/liveness-check/commit/d8f0bdfae35f3f8b324f39109555ba7bda1cf2ed))
* ignore binary ([924e8e1](https://github.com/meysam81/liveness-check/commit/924e8e1e0202ef6d3014b204918f1da3815f64da))
