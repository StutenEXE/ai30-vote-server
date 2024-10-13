## Guide d'installation

```go install github.com/StutenEXE/ai30-vote-server/cmd@latest```

Installation du projet dans le répertoire ``go/pkg/mod/github.com/'!stuten!e!x!e'`` 

## Lancement du serveur

Dans le répertoire racine du projet executer la commande suivante :

`go run cmd/launch-server.go`

Le serveur devrait être lancé et écouter sur le **port 8080**

## Fonctionnalités Implémentées

Les trois endpoints demandés (`/new_ballot`, `/vote`, `/result`) ont été implémentés avec leurs données et erreurs respectives

### Types de Scrutin

*    majority
*    approval
*    borda
*    copeland

### Erreurs Générées

*    Paramètres invalides
*    Type de vote inconnu

### Vote

*    Prise en compte de la date limite (deadline)
*    Erreur en cas de double vote
*    Non autorisé à voter
*    Scrutin inconnu

### Résultat

*    Deadline non dépassée
*    Scrutin inconnu
*    Gagnant (winner) et classement des égalités en fonction du tie-break (ranking) (notre interprétation des consignes)

### Agent

*    Création d'agents pour faciliter les requêtes
*    Tests sur les scrutins, les votes et les résultats

## Test

* Test de création et d'erreur de création de (ballot, vote, et resultat)
* Test sur les résultats d'un vote
* Pour les exectuer : `go test ./test` dans le repertoire racine du projet

#### Voir dans `` /test ``

## Architecture 

```
td5/
├─ ballot/       **(généralisation des  votes au travers d'une interface)**
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
