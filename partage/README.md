# Module Partage & API REST - Exellence Matériaux

Architecture de l'API REST Go pour le partage réseau multi-postes de la base de données.

## Structure

```
partage/
├── main.go               # Point d'entrée du serveur autonome (détection IP, arrêt propre)
├── partage_test.go       # Tests unitaires automatisés de l'API REST
├── serveur-partage.exe   # Exécutable autonome compilé
└── src/
    ├── controllers/      # Contrôleurs HTTP REST (auth, produits, factures, bons, stock, paramètres)
    ├── helpers/          # Helpers JSON, gestion des erreurs, CORS, extraction de token
    ├── models/           # Structures DTO de requêtes et réponses HTTP
    ├── routers/          # Définition des routes REST et middlewares
    └── services/         # Couche métier connectée à SQLite (db)
```

## Fonctionnalités

- **Intégration native** : Le serveur est également intégré directement dans `Exellence-Materiaux.exe` et démarre automatiquement en arrière-plan lorsque l'Administrateur se connecte sur le PC 1.
- **Port par défaut** : `8765`
- **Endpoints REST principaux** :
  - `GET  /api/health` : Statut du serveur
  - `POST /api/auth/login` : Connexion
  - `GET  /api/produits` : Catalogue des produits
  - `POST /api/produits` : Création de produit
  - `GET  /api/stock/alertes` : Alertes stock bas
  - `POST /api/stock/mouvement` : Enregistrement entrée/sortie de stock
  - `GET  /api/factures` : Historique des factures
  - `POST /api/factures` : Création et encaissement facture
  - `GET  /api/factures/:id/pdf` : Téléchargement du PDF officiel
  - `GET  /api/bons` : Historique des bons de livraison
  - `POST /api/bons` : Création d'un bon de livraison
  - `GET  /api/bons/:id/pdf` : Téléchargement du PDF officiel
  - `GET  /api/parametres` : Paramètres de l'entreprise
