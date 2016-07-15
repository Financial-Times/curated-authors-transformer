# curated-authors-transformer

[![CircleCI](https://circleci.com/gh/Financial-Times/curated-authors-transformer.svg?style=svg)](https://circleci.com/gh/Financial-Times/curated-authors-transformer)

Retrieves author data curated by editorial people and transforms it to People according to UP JSON model.
Authors data is specified by a Google spreadsheet which is accessible through [Bertha API](https://github.com/ft-interactive/bertha/wiki/Tutorial).
The service exposes endpoints for getting all the curated authors' UUIDs and for getting authors by uuid.

# How to run

## Locally: 

`go get github.com/Financial-Times/curated-authors-transformer`

`$GOPATH/bin/ ./curated-authors-transformer --bertha-source-url=<BERTHA_SOURCE_URL> --port=8080`                

```
export|set PORT=8080
export|set BERTHA_SOURCE_URL="http://.../Authors"
$GOPATH/bin/curated-authors-transformer
```

## With Docker:

`docker build -t coco/curated-authors-transformer .`

`docker run -ti --env BERTHA_SOURCE_URL=<bertha_url> coco/curated-authors-transformer`

#Endpoints

##Count
`GET /transformers/authors/__count` returns the number of available authors to be transformed as plain text.
A response example is provided below.

```
2
```

##IDs
`GET /transformers/authors/__ids` returns the list of author's UUIDs available to be transformed. 
The output is a sequence of JSON objects, however the returned `Content-Type` header is `text\plain`.
A response example is provided below.

```
{"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd2"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd5"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd9"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd8"} {"id":"5baaf5a4-2d9f-11e6-a100-1316a778acd0"} {"id":"daf5fed2-013c-468d-85c4-aee779b8aa53"} {"id":"daf5fed2-013c-468d-85c4-aee779b8aa51"} 
```

##Authors by UUID
`GET /transformers/authors/{uuid}` returns author data of the given uuid. 
A response example is provided below.

```
{
  "uuid": "daf5fed2-013c-468d-85c4-aee709b8aa53",
  "alternativeIdentifiers": {
    "TME": [
      "Q0ItMDAwMDkwMA==-QXV0auiycw=="
    ],
    "uuids": [
      "daf5fed2-013c-468d-85c4-aee779b8aa53"
    ]
  },
  "name": "Martin Wolf",
  "emailAddress": "author.email@domain.com",
  "twitterHandle": "@martinwolf_",
  "description": "Martin Wolf is chief economics commentator at the Financial Times, London. He was awarded the CBE (Commander of the British Empire) in 2000 “for services to financial journalism”",
  "descriptionXML": "<p>Martin Wolf is chief economics commentator at the Financial Times, London. He was awarded the CBE (Commander of the British Empire) in 2000 “for services to financial journalism”</p>",
  "_imageUrl": "https://example.site.com/image/martin-wolf.png"
}
```
