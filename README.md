# Exellence Matériaux

Logiciel de bureau de gestion de stock, caisse rapide et facturation pour **Exellence Matériaux**.  
Développé par **Adel Mahmoud et Sidou**.

---

## 🚀 Guide d'installation & Lancement rapide (Pour votre ami / utilisateur)

Si vous avez reçu le logiciel sous forme d'archive ZIP ou de dossier :

### Étape 1 : Extraire l'archive (Important)
Si vous avez reçu un fichier `.zip` :
1. Faites un **Clic droit** sur le fichier ZIP.
2. Cliquez sur **« Extraire tout... »** et choisissez un emplacement (par exemple sur votre **Bureau**).
3. Ouvrez le dossier extrait.
> ⚠️ **Attention** : Ne double-cliquez pas sur le logiciel directement depuis l'archive ZIP sans l'extraire, sinon Windows bloque la création de la base de données locale.

### Étape 2 : Lancer le logiciel
- Double-cliquez sur **`Exellence-Materiaux.exe`** (ou sur le script **`lancer.bat`**).

---

## 🛠️ Que faire si le logiciel ne s'ouvre pas ?

### 1. Message bleu : « Windows a protégé votre ordinateur » (SmartScreen)
C'est la sécurité par défaut de Windows lorsqu'un fichier `.exe` est téléchargé par Internet (WhatsApp, Discord, WeTransfer, Google Drive) ou copié par clé USB :
1. Cliquez sur le lien textuel **« Informations complémentaires »** situé sous le message d'alerte.
2. Le bouton **« Exécuter quand même »** apparaît : cliquez dessus. Le logiciel démarre immédiatement.
3. *Alternative* : Faites un **Clic droit sur `Exellence-Materiaux.exe`** -> **Propriétés** -> en bas de l'onglet Général, cochez la case **« Débloquer »** -> **Appliquer** -> **OK**.

### 2. L'application ne s'ouvre pas ou se ferme immédiatement (WebView2)
L'application utilise le moteur natif Microsoft WebView2 pour afficher son interface. Sur certains ordinateurs Windows 10 non mis à jour ou versions spécifiques, ce composant est absent.
- **Solution** : Téléchargez et installez l'installateur officiel Microsoft (gratuit, installation en 30 secondes) :  
  👉 [Télécharger Microsoft Edge WebView2 Evergreen Bootstrapper](https://go.microsoft.com/fwlink/p/?LinkId=2124703)

### 3. L'antivirus bloque le fichier
Certains antivirus tiers (Avast, Norton, McAfee) signalent parfois les exécutables récents non enregistrés dans leurs bases comme « Fichier inconnu ».
- **Solution** : Ajoutez une exception ou cliquez sur « Autoriser ce fichier » dans votre antivirus.

---

## 🔒 Confidentialité & Données de l'entreprise (Kbis)

Pour garantir la confidentialité totale des informations juridiques (SIREN, SIRET, RCS, adresse du siège, capital) et permettre de publier ce projet sur **GitHub en toute sécurité** :

1. Les données d'immatriculation réelles ne sont **jamais** envoyées sur GitHub.
2. Le fichier `kbis.json` contenant les informations privées est automatiquement ignoré par `.gitignore`.
3. Pour configurer les coordonnées de votre société :
   - Dupliquez le modèle `config.entreprise.example.json` et renommez-le en `kbis.json`.
   - Remplissez vos vraies informations dans `kbis.json`.
   - Au lancement, le logiciel importe automatiquement ces données dans la base locale `exellence.db`.
   - Vous pouvez également modifier ces informations à tout moment dans l'onglet **Paramètres** du logiciel.

---

## 📦 Fonctionnalités incluses

- **🛒 Caisse & Scan Rapide** : Prise en charge des douchettes code-barres USB (code + Entrée automatique), gestion du panier, calcul automatique HT/TVA/TTC, encaissement (Espèces, Carte Bancaire, Virement, Chèque).
- **📦 Catalogue Produits** : Recherche instantanée plein texte via SQLite FTS5, gestion des stocks actuels et seuils d'alerte, prix d'achat HT, vente HT et calcul automatique de la TVA.
- **📊 Gestion des Stocks** : Historique complet traçable (entrées, sorties, ventes, ajustements d'inventaire), détection et alertes visuelles de stock bas.
- **🧾 Factures & Bons de livraison** : Numérotation séquentielle certifiée sans trou (`FAC-YYYY-XXXXX`, `BL-YYYY-XXXXX`), génération et impression instantanée de PDF avec logo haute définition et mentions légales.
- **⚙️ Paramètres & Sauvegarde** : Sauvegarde instantanée de la base locale SQLite (`VACUUM INTO`) en un clic.

---

## 💻 Guide développeur (Compilation depuis les sources)

Si vous souhaitez compiler l'application vous-même depuis le code source :

### Prérequis
- [Go](https://go.dev/) >= 1.22
- [Node.js](https://nodejs.org/) >= 18 & npm
- [Wails CLI v2](https://wails.io/) : `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Commandes
1. **Mode développement (Live reload)** :
   ```bash
   wails dev
   ```

2. **Compilation de l'exécutable Windows** :
   ```bash
   wails build -tags native_webview2loader
   ```
   L'exécutable généré se trouve dans `build/bin/exellence-materiaux.exe`.

---

## 📄 Licence & Droits d'auteur

Copyright © 2026 **Adel Mahmoud et Sidou**. Tous droits réservés.
