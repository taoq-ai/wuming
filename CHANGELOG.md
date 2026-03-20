# Changelog

## [0.13.1](https://github.com/taoq-ai/wuming/compare/v0.13.0...v0.13.1) (2026-03-20)


### Documentation

* create SKILL.md for AI agent integration guidance ([#111](https://github.com/taoq-ai/wuming/issues/111)) ([959e509](https://github.com/taoq-ai/wuming/commit/959e5091d5a1224254f3a195e8d89f30ea9d0ca7)), closes [#110](https://github.com/taoq-ai/wuming/issues/110)

## [0.13.0](https://github.com/taoq-ai/wuming/compare/v0.12.0...v0.13.0) (2026-03-20)


### Features

* add allowlist/denylist for known safe values ([#104](https://github.com/taoq-ai/wuming/issues/104)) ([7dca120](https://github.com/taoq-ai/wuming/commit/7dca120b5b138af3bfc972e7a818014f144b5b5a)), closes [#100](https://github.com/taoq-ai/wuming/issues/100)
* add Australia PII detectors — TFN, Medicare, ABN, phone, postcode ([#66](https://github.com/taoq-ai/wuming/issues/66)) ([e417efa](https://github.com/taoq-ai/wuming/commit/e417efa26101f54231a9d2e1157959114eb2c3bf))
* add Brazil PII detectors — CPF, CNPJ, phone, CEP, PIS, CNH ([#70](https://github.com/taoq-ai/wuming/issues/70)) ([c4b7b1d](https://github.com/taoq-ai/wuming/commit/c4b7b1d4025e9d0feadde55f6247f877647ba44b))
* add Canada PII detectors — SIN, phone, postal code, passport ([#69](https://github.com/taoq-ai/wuming/issues/69)) ([1744bb7](https://github.com/taoq-ai/wuming/commit/1744bb7020a8b40b0bc457f4ee037c2b9f6e3f50))
* add China PII detectors ([#67](https://github.com/taoq-ai/wuming/issues/67)) ([7b805d3](https://github.com/taoq-ai/wuming/commit/7b805d309cde54f3af6c15ba60dccbe10582e05b))
* add common/global PII detectors — email, credit card, IP, URL, IBAN, MAC ([#40](https://github.com/taoq-ai/wuming/issues/40)) ([d9eb141](https://github.com/taoq-ai/wuming/commit/d9eb141e752a5984000d5e98a83701d7c0eb1d9e)), closes [#7](https://github.com/taoq-ai/wuming/issues/7)
* add compliance presets — GDPR, HIPAA, PCI-DSS, LGPD, APPI, PIPL, PIPA, DPDP, PIPEDA, Privacy Act ([#74](https://github.com/taoq-ai/wuming/issues/74)) ([dbbe4ec](https://github.com/taoq-ai/wuming/commit/dbbe4ec7b60bf09c53415cf76dd9097f63a266ee))
* add detector registry with zero-config defaults and package-level convenience functions ([#61](https://github.com/taoq-ai/wuming/issues/61)) ([7acb435](https://github.com/taoq-ai/wuming/commit/7acb4352a91fd58a21fcd1804ea2218cfbac22cf))
* add EU AI Act compliance preset for training data governance ([#81](https://github.com/taoq-ai/wuming/issues/81)) ([6e9d118](https://github.com/taoq-ai/wuming/commit/6e9d118d804c52b42cc3c2499e35aa71d4bc5a9c))
* add EU-wide PII detectors — VAT ID, passport MRZ ([#45](https://github.com/taoq-ai/wuming/issues/45)) ([b9b40ac](https://github.com/taoq-ai/wuming/commit/b9b40ac7d795a8e7b57fd48f6474d07a71a2b47a))
* add France PII detectors — NIR, NIF, phone, code postal, CNI ([#49](https://github.com/taoq-ai/wuming/issues/49)) ([888787e](https://github.com/taoq-ai/wuming/commit/888787e3271cc06a76e5789c70cddca8fcbf217f))
* add Germany PII detectors — ID card, Steuer-ID, phone, PLZ ([#50](https://github.com/taoq-ai/wuming/issues/50)) ([8a12f55](https://github.com/taoq-ai/wuming/commit/8a12f55a6e8c869b1ebe1d0daefffa94c60737b1))
* add India PII detectors ([#68](https://github.com/taoq-ai/wuming/issues/68)) ([c3ec039](https://github.com/taoq-ai/wuming/commit/c3ec039f7a36d0f18fd0f9ce00fa149a0b2f5784))
* add Japan PII detectors — My Number, Corporate Number, phone, postal, passport ([#63](https://github.com/taoq-ai/wuming/issues/63)) ([6072100](https://github.com/taoq-ai/wuming/commit/6072100884eebdb012a36b5e78a2365f07a8c868))
* add Netherlands PII detectors — BSN, phone, postal code, KvK, Dutch ID ([#48](https://github.com/taoq-ai/wuming/issues/48)) ([b243fa7](https://github.com/taoq-ai/wuming/commit/b243fa7ebf5a1b341c3aab47a98de05727508a50))
* add South Korea PII detectors — RRN, phone, postal, passport ([#64](https://github.com/taoq-ai/wuming/issues/64)) ([c7b7ec5](https://github.com/taoq-ai/wuming/commit/c7b7ec56d9f6199c65bf1c8f36749c765014713f))
* add structured data support for JSON and CSV field-level scanning ([#107](https://github.com/taoq-ai/wuming/issues/107)) ([b316fdd](https://github.com/taoq-ai/wuming/commit/b316fdd15a56d3b4a38891e7df01ff612a388859)), closes [#87](https://github.com/taoq-ai/wuming/issues/87)
* add UK PII detectors — NIN, NHS, phone, postcode, UTR ([#46](https://github.com/taoq-ai/wuming/issues/46)) ([1890580](https://github.com/taoq-ai/wuming/commit/18905803fa1d89055e859ee0337a727b4afaa8e6))
* add US PII detectors — SSN, EIN, phone, ZIP, passport, ITIN, Medicare ([#41](https://github.com/taoq-ai/wuming/issues/41)) ([cc8ac96](https://github.com/taoq-ai/wuming/commit/cc8ac96dfa3c89913e11aca43c65297c59c55082)), closes [#8](https://github.com/taoq-ai/wuming/issues/8)
* consistent redaction — same PII value maps to same replacement ([#106](https://github.com/taoq-ai/wuming/issues/106)) ([94e869c](https://github.com/taoq-ai/wuming/commit/94e869c2fcf1ea98522eebf694f5c5d0ba403749)), closes [#88](https://github.com/taoq-ai/wuming/issues/88)
* v0.1.0 foundation — hexagonal architecture scaffold ([#30](https://github.com/taoq-ai/wuming/issues/30)) ([ea52c0d](https://github.com/taoq-ai/wuming/commit/ea52c0df217fec5653056b99bba94691e9dd8925))


### Bug Fixes

* force first release as v0.1.0 via release-as ([#37](https://github.com/taoq-ai/wuming/issues/37)) ([8e76b7e](https://github.com/taoq-ai/wuming/commit/8e76b7ea82f0c3dc89d2da683d4a827fc115e9a7)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* improve benchmark PR comment formatting ([#85](https://github.com/taoq-ai/wuming/issues/85)) ([58a5ce6](https://github.com/taoq-ai/wuming/commit/58a5ce697279e58a37ba33138727be8e6c463486))
* remove bump-patch-for-minor-pre-major so feat: bumps minor ([#43](https://github.com/taoq-ai/wuming/issues/43)) ([3322dc8](https://github.com/taoq-ai/wuming/commit/3322dc812481fa9ab4146a79f86e0a1394c8a292))
* set initial manifest version to 0.0.1 for correct pre-1.0 bumping ([#34](https://github.com/taoq-ai/wuming/issues/34)) ([a39fcf0](https://github.com/taoq-ai/wuming/commit/a39fcf0ac7b23a4f092d76ab120c90382e52d8fd)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* use fence_code_format for Mermaid in mkdocs.yml ([#58](https://github.com/taoq-ai/wuming/issues/58)) ([b38b3a9](https://github.com/taoq-ai/wuming/commit/b38b3a9eed1f9090ee7131e8a58426ae5d4303e6))
* use PAT for release-please to allow PR creation ([#32](https://github.com/taoq-ai/wuming/issues/32)) ([6d4b933](https://github.com/taoq-ai/wuming/commit/6d4b9330efe6f6fdf2c8aa45594ad0cffea3635e)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* WithPreset returns error on unknown name instead of panicking ([#105](https://github.com/taoq-ai/wuming/issues/105)) ([0b8c8fd](https://github.com/taoq-ai/wuming/commit/0b8c8fd7e5a59f5c31d3b6c7e3c275de4daa6632))


### Documentation

* add CITATION.cff for academic citation ([#51](https://github.com/taoq-ai/wuming/issues/51)) ([c6ad237](https://github.com/taoq-ai/wuming/commit/c6ad237a93f2d5a7fc177ccbc11cd0d650e0045a))
* add MkDocs site with Material theme and GitHub Pages deployment ([#55](https://github.com/taoq-ai/wuming/issues/55)) ([0bd2a41](https://github.com/taoq-ai/wuming/commit/0bd2a41a9e9e36d10fb021fab5f4faae33a48faa))
* add README, CONTRIBUTING, ARCHITECTURE, SECURITY, and examples ([#53](https://github.com/taoq-ai/wuming/issues/53)) ([87e37bc](https://github.com/taoq-ai/wuming/commit/87e37bc74d898fd8588e7368d6493e3a5c5fffd2))
* update README and docs for all 14 locales, presets, and zero-config API ([#78](https://github.com/taoq-ai/wuming/issues/78)) ([9dc9d22](https://github.com/taoq-ai/wuming/commit/9dc9d2261f7e3289fd66111770de85dfe664a16a))


### Miscellaneous

* **main:** release 0.1.0 ([#38](https://github.com/taoq-ai/wuming/issues/38)) ([cbb5e8b](https://github.com/taoq-ai/wuming/commit/cbb5e8bdddeabc427871fda0d2fd5dfda51e911c))
* **main:** release 0.10.0 ([#82](https://github.com/taoq-ai/wuming/issues/82)) ([4278fca](https://github.com/taoq-ai/wuming/commit/4278fca1f2384764486407ae16297bf38f6b99d6))
* **main:** release 0.10.1 ([#103](https://github.com/taoq-ai/wuming/issues/103)) ([d58a7e1](https://github.com/taoq-ai/wuming/commit/d58a7e1abce2909ebf49e5336ff531c39b88b635))
* **main:** release 0.11.0 ([#108](https://github.com/taoq-ai/wuming/issues/108)) ([58e7ec9](https://github.com/taoq-ai/wuming/commit/58e7ec9c3a2f4e13af065e7bc2d010383761befb))
* **main:** release 0.11.1 ([#109](https://github.com/taoq-ai/wuming/issues/109)) ([c7b0de9](https://github.com/taoq-ai/wuming/commit/c7b0de95908ca5f42329f39c57545dc1c2317a02))
* **main:** release 0.12.0 ([#112](https://github.com/taoq-ai/wuming/issues/112)) ([f988116](https://github.com/taoq-ai/wuming/commit/f9881162e9971ed3cc17d0ba383a873927e67a73))
* **main:** release 0.2.0 ([#44](https://github.com/taoq-ai/wuming/issues/44)) ([699851d](https://github.com/taoq-ai/wuming/commit/699851d107b41e974e08530274b47c8f935347af))
* **main:** release 0.3.0 ([#47](https://github.com/taoq-ai/wuming/issues/47)) ([9354e57](https://github.com/taoq-ai/wuming/commit/9354e5757ea33c37a555aaba7266bcf4560d3d16))
* **main:** release 0.3.1 ([#52](https://github.com/taoq-ai/wuming/issues/52)) ([6a60c56](https://github.com/taoq-ai/wuming/commit/6a60c56b9cc6e87f9b16fcd2e3b96248ea682d68))
* **main:** release 0.3.2 ([#54](https://github.com/taoq-ai/wuming/issues/54)) ([f37769b](https://github.com/taoq-ai/wuming/commit/f37769bd2b43a6f89a11751f5b4c8e9977babb8e))
* **main:** release 0.3.3 ([#56](https://github.com/taoq-ai/wuming/issues/56)) ([358ecfb](https://github.com/taoq-ai/wuming/commit/358ecfbb3b85f385855c990161eedfee5fc0a23c))
* **main:** release 0.3.4 ([#59](https://github.com/taoq-ai/wuming/issues/59)) ([b4d604f](https://github.com/taoq-ai/wuming/commit/b4d604fe0000b2e71b3ef4569263c089753d19c3))
* **main:** release 0.4.0 ([#60](https://github.com/taoq-ai/wuming/issues/60)) ([430b043](https://github.com/taoq-ai/wuming/commit/430b0435d7a7958075f7fc7697c52cfaed7ca581))
* **main:** release 0.5.0 ([#62](https://github.com/taoq-ai/wuming/issues/62)) ([8856379](https://github.com/taoq-ai/wuming/commit/8856379b8121607d2ca64e95dc7be77c9b385916))
* **main:** release 0.6.0 ([#65](https://github.com/taoq-ai/wuming/issues/65)) ([6906afa](https://github.com/taoq-ai/wuming/commit/6906afa5d1e9c77e49b44a56cf6cb1eb8fe65183))
* **main:** release 0.7.0 ([#71](https://github.com/taoq-ai/wuming/issues/71)) ([ba52ec0](https://github.com/taoq-ai/wuming/commit/ba52ec04a6897d3bf87ae2c60e78ccca0a9d4b7d))
* **main:** release 0.8.0 ([#72](https://github.com/taoq-ai/wuming/issues/72)) ([8142416](https://github.com/taoq-ai/wuming/commit/81424169bea543e21408f31af88003e7e383e186))
* **main:** release 0.9.0 ([#75](https://github.com/taoq-ai/wuming/issues/75)) ([a5b4bba](https://github.com/taoq-ai/wuming/commit/a5b4bba206e9ae3f45b7533a9fdc801e2faf0bcc))
* **main:** release 0.9.1 ([#79](https://github.com/taoq-ai/wuming/issues/79)) ([880844e](https://github.com/taoq-ai/wuming/commit/880844e80bcce5928634d4afc9d462c4f3330220))
* remove release-as now that v0.1.0 is published ([#39](https://github.com/taoq-ai/wuming/issues/39)) ([c51aa5a](https://github.com/taoq-ai/wuming/commit/c51aa5afb2e1132c25ac71571a24b9a0b4b6561f))

## [0.12.0](https://github.com/taoq-ai/wuming/compare/v0.11.1...v0.12.0) (2026-03-20)


### Features

* consistent redaction — same PII value maps to same replacement ([#106](https://github.com/taoq-ai/wuming/issues/106)) ([94e869c](https://github.com/taoq-ai/wuming/commit/94e869c2fcf1ea98522eebf694f5c5d0ba403749)), closes [#88](https://github.com/taoq-ai/wuming/issues/88)

## [0.11.1](https://github.com/taoq-ai/wuming/compare/v0.11.0...v0.11.1) (2026-03-20)


### Bug Fixes

* WithPreset returns error on unknown name instead of panicking ([#105](https://github.com/taoq-ai/wuming/issues/105)) ([0b8c8fd](https://github.com/taoq-ai/wuming/commit/0b8c8fd7e5a59f5c31d3b6c7e3c275de4daa6632))

## [0.11.0](https://github.com/taoq-ai/wuming/compare/v0.10.1...v0.11.0) (2026-03-20)


### Features

* add allowlist/denylist for known safe values ([#104](https://github.com/taoq-ai/wuming/issues/104)) ([7dca120](https://github.com/taoq-ai/wuming/commit/7dca120b5b138af3bfc972e7a818014f144b5b5a)), closes [#100](https://github.com/taoq-ai/wuming/issues/100)
* add structured data support for JSON and CSV field-level scanning ([#107](https://github.com/taoq-ai/wuming/issues/107)) ([b316fdd](https://github.com/taoq-ai/wuming/commit/b316fdd15a56d3b4a38891e7df01ff612a388859)), closes [#87](https://github.com/taoq-ai/wuming/issues/87)

## [0.10.1](https://github.com/taoq-ai/wuming/compare/v0.10.0...v0.10.1) (2026-03-20)


### Bug Fixes

* improve benchmark PR comment formatting ([#85](https://github.com/taoq-ai/wuming/issues/85)) ([58a5ce6](https://github.com/taoq-ai/wuming/commit/58a5ce697279e58a37ba33138727be8e6c463486))

## [0.10.0](https://github.com/taoq-ai/wuming/compare/v0.9.1...v0.10.0) (2026-03-19)


### Features

* add EU AI Act compliance preset for training data governance ([#81](https://github.com/taoq-ai/wuming/issues/81)) ([6e9d118](https://github.com/taoq-ai/wuming/commit/6e9d118d804c52b42cc3c2499e35aa71d4bc5a9c))

## [0.9.1](https://github.com/taoq-ai/wuming/compare/v0.9.0...v0.9.1) (2026-03-19)


### Documentation

* update README and docs for all 14 locales, presets, and zero-config API ([#78](https://github.com/taoq-ai/wuming/issues/78)) ([9dc9d22](https://github.com/taoq-ai/wuming/commit/9dc9d2261f7e3289fd66111770de85dfe664a16a))

## [0.9.0](https://github.com/taoq-ai/wuming/compare/v0.8.0...v0.9.0) (2026-03-19)


### Features

* add compliance presets — GDPR, HIPAA, PCI-DSS, LGPD, APPI, PIPL, PIPA, DPDP, PIPEDA, Privacy Act ([#74](https://github.com/taoq-ai/wuming/issues/74)) ([dbbe4ec](https://github.com/taoq-ai/wuming/commit/dbbe4ec7b60bf09c53415cf76dd9097f63a266ee))

## [0.8.0](https://github.com/taoq-ai/wuming/compare/v0.7.0...v0.8.0) (2026-03-19)


### Features

* add Brazil PII detectors — CPF, CNPJ, phone, CEP, PIS, CNH ([#70](https://github.com/taoq-ai/wuming/issues/70)) ([c4b7b1d](https://github.com/taoq-ai/wuming/commit/c4b7b1d4025e9d0feadde55f6247f877647ba44b))

## [0.7.0](https://github.com/taoq-ai/wuming/compare/v0.6.0...v0.7.0) (2026-03-19)


### Features

* add Australia PII detectors — TFN, Medicare, ABN, phone, postcode ([#66](https://github.com/taoq-ai/wuming/issues/66)) ([e417efa](https://github.com/taoq-ai/wuming/commit/e417efa26101f54231a9d2e1157959114eb2c3bf))
* add Canada PII detectors — SIN, phone, postal code, passport ([#69](https://github.com/taoq-ai/wuming/issues/69)) ([1744bb7](https://github.com/taoq-ai/wuming/commit/1744bb7020a8b40b0bc457f4ee037c2b9f6e3f50))
* add China PII detectors ([#67](https://github.com/taoq-ai/wuming/issues/67)) ([7b805d3](https://github.com/taoq-ai/wuming/commit/7b805d309cde54f3af6c15ba60dccbe10582e05b))
* add India PII detectors ([#68](https://github.com/taoq-ai/wuming/issues/68)) ([c3ec039](https://github.com/taoq-ai/wuming/commit/c3ec039f7a36d0f18fd0f9ce00fa149a0b2f5784))
* add South Korea PII detectors — RRN, phone, postal, passport ([#64](https://github.com/taoq-ai/wuming/issues/64)) ([c7b7ec5](https://github.com/taoq-ai/wuming/commit/c7b7ec56d9f6199c65bf1c8f36749c765014713f))

## [0.6.0](https://github.com/taoq-ai/wuming/compare/v0.5.0...v0.6.0) (2026-03-19)


### Features

* add Japan PII detectors — My Number, Corporate Number, phone, postal, passport ([#63](https://github.com/taoq-ai/wuming/issues/63)) ([6072100](https://github.com/taoq-ai/wuming/commit/6072100884eebdb012a36b5e78a2365f07a8c868))

## [0.5.0](https://github.com/taoq-ai/wuming/compare/v0.4.0...v0.5.0) (2026-03-19)


### Features

* add detector registry with zero-config defaults and package-level convenience functions ([#61](https://github.com/taoq-ai/wuming/issues/61)) ([7acb435](https://github.com/taoq-ai/wuming/commit/7acb4352a91fd58a21fcd1804ea2218cfbac22cf))

## [0.4.0](https://github.com/taoq-ai/wuming/compare/v0.3.4...v0.4.0) (2026-03-19)


### Features

* add common/global PII detectors — email, credit card, IP, URL, IBAN, MAC ([#40](https://github.com/taoq-ai/wuming/issues/40)) ([d9eb141](https://github.com/taoq-ai/wuming/commit/d9eb141e752a5984000d5e98a83701d7c0eb1d9e)), closes [#7](https://github.com/taoq-ai/wuming/issues/7)
* add EU-wide PII detectors — VAT ID, passport MRZ ([#45](https://github.com/taoq-ai/wuming/issues/45)) ([b9b40ac](https://github.com/taoq-ai/wuming/commit/b9b40ac7d795a8e7b57fd48f6474d07a71a2b47a))
* add France PII detectors — NIR, NIF, phone, code postal, CNI ([#49](https://github.com/taoq-ai/wuming/issues/49)) ([888787e](https://github.com/taoq-ai/wuming/commit/888787e3271cc06a76e5789c70cddca8fcbf217f))
* add Germany PII detectors — ID card, Steuer-ID, phone, PLZ ([#50](https://github.com/taoq-ai/wuming/issues/50)) ([8a12f55](https://github.com/taoq-ai/wuming/commit/8a12f55a6e8c869b1ebe1d0daefffa94c60737b1))
* add Netherlands PII detectors — BSN, phone, postal code, KvK, Dutch ID ([#48](https://github.com/taoq-ai/wuming/issues/48)) ([b243fa7](https://github.com/taoq-ai/wuming/commit/b243fa7ebf5a1b341c3aab47a98de05727508a50))
* add UK PII detectors — NIN, NHS, phone, postcode, UTR ([#46](https://github.com/taoq-ai/wuming/issues/46)) ([1890580](https://github.com/taoq-ai/wuming/commit/18905803fa1d89055e859ee0337a727b4afaa8e6))
* add US PII detectors — SSN, EIN, phone, ZIP, passport, ITIN, Medicare ([#41](https://github.com/taoq-ai/wuming/issues/41)) ([cc8ac96](https://github.com/taoq-ai/wuming/commit/cc8ac96dfa3c89913e11aca43c65297c59c55082)), closes [#8](https://github.com/taoq-ai/wuming/issues/8)
* v0.1.0 foundation — hexagonal architecture scaffold ([#30](https://github.com/taoq-ai/wuming/issues/30)) ([ea52c0d](https://github.com/taoq-ai/wuming/commit/ea52c0df217fec5653056b99bba94691e9dd8925))


### Bug Fixes

* force first release as v0.1.0 via release-as ([#37](https://github.com/taoq-ai/wuming/issues/37)) ([8e76b7e](https://github.com/taoq-ai/wuming/commit/8e76b7ea82f0c3dc89d2da683d4a827fc115e9a7)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* remove bump-patch-for-minor-pre-major so feat: bumps minor ([#43](https://github.com/taoq-ai/wuming/issues/43)) ([3322dc8](https://github.com/taoq-ai/wuming/commit/3322dc812481fa9ab4146a79f86e0a1394c8a292))
* set initial manifest version to 0.0.1 for correct pre-1.0 bumping ([#34](https://github.com/taoq-ai/wuming/issues/34)) ([a39fcf0](https://github.com/taoq-ai/wuming/commit/a39fcf0ac7b23a4f092d76ab120c90382e52d8fd)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* use fence_code_format for Mermaid in mkdocs.yml ([#58](https://github.com/taoq-ai/wuming/issues/58)) ([b38b3a9](https://github.com/taoq-ai/wuming/commit/b38b3a9eed1f9090ee7131e8a58426ae5d4303e6))
* use PAT for release-please to allow PR creation ([#32](https://github.com/taoq-ai/wuming/issues/32)) ([6d4b933](https://github.com/taoq-ai/wuming/commit/6d4b9330efe6f6fdf2c8aa45594ad0cffea3635e)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)


### Documentation

* add CITATION.cff for academic citation ([#51](https://github.com/taoq-ai/wuming/issues/51)) ([c6ad237](https://github.com/taoq-ai/wuming/commit/c6ad237a93f2d5a7fc177ccbc11cd0d650e0045a))
* add MkDocs site with Material theme and GitHub Pages deployment ([#55](https://github.com/taoq-ai/wuming/issues/55)) ([0bd2a41](https://github.com/taoq-ai/wuming/commit/0bd2a41a9e9e36d10fb021fab5f4faae33a48faa))
* add README, CONTRIBUTING, ARCHITECTURE, SECURITY, and examples ([#53](https://github.com/taoq-ai/wuming/issues/53)) ([87e37bc](https://github.com/taoq-ai/wuming/commit/87e37bc74d898fd8588e7368d6493e3a5c5fffd2))


### Miscellaneous

* **main:** release 0.1.0 ([#38](https://github.com/taoq-ai/wuming/issues/38)) ([cbb5e8b](https://github.com/taoq-ai/wuming/commit/cbb5e8bdddeabc427871fda0d2fd5dfda51e911c))
* **main:** release 0.2.0 ([#44](https://github.com/taoq-ai/wuming/issues/44)) ([699851d](https://github.com/taoq-ai/wuming/commit/699851d107b41e974e08530274b47c8f935347af))
* **main:** release 0.3.0 ([#47](https://github.com/taoq-ai/wuming/issues/47)) ([9354e57](https://github.com/taoq-ai/wuming/commit/9354e5757ea33c37a555aaba7266bcf4560d3d16))
* **main:** release 0.3.1 ([#52](https://github.com/taoq-ai/wuming/issues/52)) ([6a60c56](https://github.com/taoq-ai/wuming/commit/6a60c56b9cc6e87f9b16fcd2e3b96248ea682d68))
* **main:** release 0.3.2 ([#54](https://github.com/taoq-ai/wuming/issues/54)) ([f37769b](https://github.com/taoq-ai/wuming/commit/f37769bd2b43a6f89a11751f5b4c8e9977babb8e))
* **main:** release 0.3.3 ([#56](https://github.com/taoq-ai/wuming/issues/56)) ([358ecfb](https://github.com/taoq-ai/wuming/commit/358ecfbb3b85f385855c990161eedfee5fc0a23c))
* **main:** release 0.3.4 ([#59](https://github.com/taoq-ai/wuming/issues/59)) ([b4d604f](https://github.com/taoq-ai/wuming/commit/b4d604fe0000b2e71b3ef4569263c089753d19c3))
* remove release-as now that v0.1.0 is published ([#39](https://github.com/taoq-ai/wuming/issues/39)) ([c51aa5a](https://github.com/taoq-ai/wuming/commit/c51aa5afb2e1132c25ac71571a24b9a0b4b6561f))

## [0.3.4](https://github.com/taoq-ai/wuming/compare/v0.3.3...v0.3.4) (2026-03-19)


### Bug Fixes

* use fence_code_format for Mermaid in mkdocs.yml ([#58](https://github.com/taoq-ai/wuming/issues/58)) ([b38b3a9](https://github.com/taoq-ai/wuming/commit/b38b3a9eed1f9090ee7131e8a58426ae5d4303e6))

## [0.3.3](https://github.com/taoq-ai/wuming/compare/v0.3.2...v0.3.3) (2026-03-19)


### Documentation

* add MkDocs site with Material theme and GitHub Pages deployment ([#55](https://github.com/taoq-ai/wuming/issues/55)) ([0bd2a41](https://github.com/taoq-ai/wuming/commit/0bd2a41a9e9e36d10fb021fab5f4faae33a48faa))

## [0.3.2](https://github.com/taoq-ai/wuming/compare/v0.3.1...v0.3.2) (2026-03-19)


### Documentation

* add README, CONTRIBUTING, ARCHITECTURE, SECURITY, and examples ([#53](https://github.com/taoq-ai/wuming/issues/53)) ([87e37bc](https://github.com/taoq-ai/wuming/commit/87e37bc74d898fd8588e7368d6493e3a5c5fffd2))

## [0.3.1](https://github.com/taoq-ai/wuming/compare/v0.3.0...v0.3.1) (2026-03-19)


### Documentation

* add CITATION.cff for academic citation ([#51](https://github.com/taoq-ai/wuming/issues/51)) ([c6ad237](https://github.com/taoq-ai/wuming/commit/c6ad237a93f2d5a7fc177ccbc11cd0d650e0045a))

## [0.3.0](https://github.com/taoq-ai/wuming/compare/v0.2.0...v0.3.0) (2026-03-19)


### Features

* add EU-wide PII detectors — VAT ID, passport MRZ ([#45](https://github.com/taoq-ai/wuming/issues/45)) ([b9b40ac](https://github.com/taoq-ai/wuming/commit/b9b40ac7d795a8e7b57fd48f6474d07a71a2b47a))
* add France PII detectors — NIR, NIF, phone, code postal, CNI ([#49](https://github.com/taoq-ai/wuming/issues/49)) ([888787e](https://github.com/taoq-ai/wuming/commit/888787e3271cc06a76e5789c70cddca8fcbf217f))
* add Germany PII detectors — ID card, Steuer-ID, phone, PLZ ([#50](https://github.com/taoq-ai/wuming/issues/50)) ([8a12f55](https://github.com/taoq-ai/wuming/commit/8a12f55a6e8c869b1ebe1d0daefffa94c60737b1))
* add Netherlands PII detectors — BSN, phone, postal code, KvK, Dutch ID ([#48](https://github.com/taoq-ai/wuming/issues/48)) ([b243fa7](https://github.com/taoq-ai/wuming/commit/b243fa7ebf5a1b341c3aab47a98de05727508a50))
* add UK PII detectors — NIN, NHS, phone, postcode, UTR ([#46](https://github.com/taoq-ai/wuming/issues/46)) ([1890580](https://github.com/taoq-ai/wuming/commit/18905803fa1d89055e859ee0337a727b4afaa8e6))

## [0.2.0](https://github.com/taoq-ai/wuming/compare/v0.1.0...v0.2.0) (2026-03-19)


### Features

* add common/global PII detectors — email, credit card, IP, URL, IBAN, MAC ([#40](https://github.com/taoq-ai/wuming/issues/40)) ([d9eb141](https://github.com/taoq-ai/wuming/commit/d9eb141e752a5984000d5e98a83701d7c0eb1d9e)), closes [#7](https://github.com/taoq-ai/wuming/issues/7)
* add US PII detectors — SSN, EIN, phone, ZIP, passport, ITIN, Medicare ([#41](https://github.com/taoq-ai/wuming/issues/41)) ([cc8ac96](https://github.com/taoq-ai/wuming/commit/cc8ac96dfa3c89913e11aca43c65297c59c55082)), closes [#8](https://github.com/taoq-ai/wuming/issues/8)


### Bug Fixes

* remove bump-patch-for-minor-pre-major so feat: bumps minor ([#43](https://github.com/taoq-ai/wuming/issues/43)) ([3322dc8](https://github.com/taoq-ai/wuming/commit/3322dc812481fa9ab4146a79f86e0a1394c8a292))


### Miscellaneous

* remove release-as now that v0.1.0 is published ([#39](https://github.com/taoq-ai/wuming/issues/39)) ([c51aa5a](https://github.com/taoq-ai/wuming/commit/c51aa5afb2e1132c25ac71571a24b9a0b4b6561f))

## 0.1.0 (2026-03-19)


### Features

* v0.1.0 foundation — hexagonal architecture scaffold ([#30](https://github.com/taoq-ai/wuming/issues/30)) ([ea52c0d](https://github.com/taoq-ai/wuming/commit/ea52c0df217fec5653056b99bba94691e9dd8925))


### Bug Fixes

* force first release as v0.1.0 via release-as ([#37](https://github.com/taoq-ai/wuming/issues/37)) ([8e76b7e](https://github.com/taoq-ai/wuming/commit/8e76b7ea82f0c3dc89d2da683d4a827fc115e9a7)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* set initial manifest version to 0.0.1 for correct pre-1.0 bumping ([#34](https://github.com/taoq-ai/wuming/issues/34)) ([a39fcf0](https://github.com/taoq-ai/wuming/commit/a39fcf0ac7b23a4f092d76ab120c90382e52d8fd)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
* use PAT for release-please to allow PR creation ([#32](https://github.com/taoq-ai/wuming/issues/32)) ([6d4b933](https://github.com/taoq-ai/wuming/commit/6d4b9330efe6f6fdf2c8aa45594ad0cffea3635e)), closes [#25](https://github.com/taoq-ai/wuming/issues/25)
