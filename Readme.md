# Auth-Service

## ✨ Description
Le microservice `auth-service` gère l'ensemble du cycle d'authentification et d'autorisation de l'application Xpensify. Il est conçu pour être évolutif, sécurisé, et interopérable avec les autres microservices via un système de publication/subscription (pub/sub) basé sur Redis.

## 🛡️ Fonctionnalités principales
- Authentification via JWT (login/register/refresh/logout)
- Récupération de l'utilisateur connecté (`/me`)
- Gestion des services : CRUD sur les services utilisables par les utilisateurs
- Gestion des permissions : assignation/retrait de permissions aux utilisateurs
- Gestion des utilisateurs et leurs permissions (modification et listing)
- Publication d'événements Redis lors des modifications de permissions

## 🚩 Stack technique
- **Langage** : Golang
- **Framework** : [Fiber](https://gofiber.io/) (rapide, inspiré d'Express.js)
- **Base de données** : PostgreSQL (via `gorm`)
- **Authentification** : JWT avec vérification via `middleware`
- **Pub/Sub** : Redis (pour synchronisation inter-microservices)
- **Conteneurisation** : Docker (et préparation pour Kubernetes)

## 📂 Arborescence
```
.
├── database/                  # Connexion à PostgreSQL
├── devops/                   # Docker et Kubernetes
│   ├── compose/
│   └── k8s/
├── features/                 # Handlers HTTP (sync et async à venir)
├── models/                   # Modèles de données (User, Permission, etc.)
├── pubsub/                   # Gestion Redis Pub/Sub
├── routes/                   # Définition des routes
├── security/                 # JWT, middlewares de rôles, mot de passe
├── services/                 # Logique métier (UserService, AuthService...)
├── main.go                   # Point d'entrée
├── Dockerfile
├── go.mod / go.sum
```

## 🔄 Endpoints disponibles
Prefixe : `/api/auth`

### ✅ Authentification
- `POST /register` : Créer un compte
- `POST /login` : Se connecter
- `POST /refresh` : Régénérer un token JWT
- `POST /logout` : Invalider la session (via Redis)
- `GET /me` : Infos de l'utilisateur connecté (protégé)

### 📑 Services (requiert `superuser`)
- `GET /services`
- `POST /services`
- `PUT /services/:id`
- `DELETE /services/:id`

### 🔒 Permissions (requiert `superuser`)
- `POST /permissions`
- `DELETE /permissions/:id`

### 🤵🏻 Utilisateurs (requiert `superuser`)
- `GET /users`
- `PUT /users/:id/permissions`

## 🧰 Redis Pub/Sub
Lorsqu'une permission est modifiée, le service publie un message Redis sur un channel dédié. Ce message peut être écouté par d'autres services pour mettre à jour leurs droits en local (cache, JWT validation, etc.).

## 🚚 Déploiement & DevOps
- Dockerfile fourni pour l'image de base
- Docker Compose pour environnement local
- YAML Kubernetes pour PostgreSQL présent dans `devops/k8s`
- CI/CD prévu via Azure Pipelines (ou GitHub Actions)

## 📌 Prochaines étapes
- Mise en place du CI/CD
- Ajout de tests unitaires et d'intégration
- Déploiement Kubernetes global
- Intégration OAuth Google/Facebook
- Gestion fine des scopes/claims dans JWT

---

Ce `auth-service` est au cœur de l'architecture microservices de Xpensify, garantissant l'identité, la sécurité et la coordination des accès pour chaque utilisateur et service.

