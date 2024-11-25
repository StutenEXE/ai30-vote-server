# Serveur de Vote - TD5 AI30

|Information|Valeur|
|-|-|
|Auteurs|Martins Clément|
||Bidaux Alexandre|
|Version|1.0.0|

## Guide d'installation

Installation à partir de la commande suivante : 

```go install github.com/StutenEXE/ai30-vote-server/cmd@latest```

Installation du projet dans le répertoire ``go/pkg/mod/github.com/'!stuten!e!x!e'`` 

## Lancement du serveur

Dans le répertoire racine du projet executer la commande suivante :

`go run cmd/launch-server.go`

Le serveur devrait être lancé et écouter sur le **port 8080**

## Fonctionnalités implémentées

Les trois endpoints demandés (`/new_ballot`, `/vote`, `/result`) ont été implémentés avec leurs données et erreurs respectives

* Types de Scrutin : majority, approval, borda, copeland
* Erreurs traitées : toutes les erreurs demandées sont renvoyées avec le code d'erreur correspondant et parfois un message decrivant l'erreur 
    * Création de ballot
        * Paramètres invalides
        * Type de vote inconnu
    * Vote
        * Prise en compte de la date limite (deadline)
        * Erreur en cas de double vote
        * Non autorisé à voter
        * Scrutin inconnu
    * Résultat
        * Deadline non dépassée
        * Scrutin inconnu
        * Gagnant (winner) et classement des égalités en fonction du tie-break (ranking) (notre interprétation des consignes)
* Agent
   * Création d'agents pour faciliter l'envoi de requêtes
   * Tests sur les scrutins, les votes et les résultats

## Tests

* Tests d'erreurs de création des ballots, votes, et resultats
* Tests sur les résultats d'un vote
* Pour exectuer les tests : lancer le serveur puis les test (`go run ./cmd/launch-server.go` et `go test ./test` dans le repertoire racine du projet)

## Architecture 

```
td5/
├─ ballot/       **(généralisation des votes au travers d'une interface)**
│  ├─ approvalballot.go
│  ├─ ballot.go  **(interface)**
│  ├─ bordaballot.go
│  ├─ copelandballot.go
│  ├─ majorityballot.go
├─ cmd/          **(codes de lancement de l'application)**
│  ├─ launch-server.go
├─ comsoc/       **(logiques de vote, provient du td4)**
├─ test/
│  ├─ ballot_test.go
│  ├─ result_test.go
│  ├─ vote_test.go
├─ voteserveragent/
│  ├─ server.go
├─ voteclientagent/
│  ├─ client.go
├─ types.go      **(utilisés partout dans le code)**
```

## API de notre web service

### Endpoint `/new_ballot`

* Requête : `POST`
* Types de scrutins implémentés : majority, approval, borda, copeland
* Objet `JSON` envoyé

| Propriété | Type         | Exemple de valeurs possibles                  |
|-----------|--------------|-----------------------------------------------|
| rule      | string       | "majority", "borda", "approval", "copeland"   |
| deadline  | string       | "2024-10-09T14:35:19+02:00" (format RFC 3339) |
| voter-ids | [string,...] | ["ag_id1", "ag_id2", "ag_id3"]                |
| #alts     | int          | 12                                            |
| tie-break | [int,...]    | [4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10]       |

* Codes retour

| Code retour | Signification           |
|-------------|-------------------------|
| 201         | vote créé               |
| 400         | bad request             |
| 501         | not implemented         |

> Remarque : un message peut être renvoyé avec le code de retour pour décrire l'erreur

* Objet `JSON` renvoyé (si `201`)

| Propriété | Type   | Exemple de valeurs possibles |
|-----------|--------|------------------------------|
| ballot-id | string | "scrutin12"                  |

### Endpoint `/vote`

* Requête : `POST`
* Objet `JSON` envoyé

| Propriété | Type       | Exemple de valeurs possibles |
|-----------|------------|------------------------------|
| agent-id  | string     | "ag_id1"                     |
| ballot-id | string     | "scrutin12"                  |
| prefs     | [int,...]  | [1, 2, 4, 3]                 |
| options   | [int,...]  | [3]                          |

> Remarque : `options` est facultatif et permet de passer des renseignements supplémentaires (par exemple le seuil d'acceptation en approval)

* Codes retour

| Code retour | Signification           |
|-------------|-------------------------|
| 200         | vote pris en compte     |
| 400         | bad request             |
| 403         | vote déjà effectué      |
| 503         | la deadline est dépassée|

> Remarque : un message peut être renvoyé avec le code de retour pour décrire l'erreur

### Endpoint `/result`

* Requête : `POST`
* Objet `JSON` envoyé

| Propriété | Type   | Exemple de valeurs possibles |
|-----------|--------|------------------------------|
| ballot-id | string | "scrutin12"                  |

* Code retour

| Code retour | Signification |
|-------------|---------------|
| 200         | OK            |
| 425         | Too early     |
| 404         | Not Found     |

> Remarque : un message peut être renvoyé avec le code de retour pour décrire l'erreur


* Objet `JSON` renvoyé (si `200`)

| Propriété | Type       | Exemple de valeurs possibles |
|-----------|------------|------------------------------|
| winner    | int        | 4                            |
| ranking   | [int,...]  | [2, 1, 4, 3]                 |

> Remarque : tel que nous l'avons compris, `ranking` renvoie les gagnants de l'élection sans faire de tie-break
