# SwimTools

This repo contains a golang service that provides a set of useful tools to operate with swimming exchange files. It is based on GIN and offers REST endpoints for all features.

## 📋 DSV Parser

The [DSV Parser](https://github.com/konrad2002/dsvparser) can be used with this tool.

## 📋 LENEX Parser

The [LENEX Parser](https://github.com/konrad2002/lenexparser) can be used with this tool.

## DSV <-> LENEX Converter

The converter can be used to convert DSV files to LENEX and vice versa. The following table shows the current support for the different file types:

| **DSV File Type**         | **DSV -> LENEX** | **LENEX -> DSV** |
|---------------------------|------------------|------------------|
| Wettkampfdefinitionsliste | 🟥               | 🟥               |
| Vereinsmeldeliste         | 🟥               | 🟥               |
| Wettkampfergebnisliste    | 🟧               | 🟧               |
| Vereinsergebnisliste      | 🟥               | 🟥               | 

- 🟩 - supported
- 🟨 - partially supported
- 🟧 - planned
- 🟥 - not supported

