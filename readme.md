## Installation Guide

```go install github.com/StutenEXE/ai30-vote-server/cmd@latest```

## Fonctionnalités Implémentées
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
*    Gagnant et éventuellement un classement

### Agent

*    Création d'agents pour faciliter les requêtes
*    Tests sur les scrutins, les votes et les résultats

## Test
* Test de création et d'erreur de création de (ballot, vote, et resultat)
* Test sur les résultats d'un vote

#### Voir dans `` /test ``
