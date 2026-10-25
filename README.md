# Gestion de Stock

Logiciel de bureau tout-en-un : Caisse & scan rapide, gestion des stocks, facturation et bons de livraison.  
Conçu et développé par **Adel Mahmoud et Sidou**.

---

## ⚡ Prérequis (Pour l'utilisateur)

**AUCUNE CONNAISSANCE EN CODE N'EST NÉCESSAIRE.**  
Le logiciel est **déjà compilé et 100% prêt à l'emploi**.

- ❌ **Pas besoin** d'installer Go.
- ❌ **Pas besoin** d'installer Node.js ni de terminal.
- ❌ **Pas de lignes de commande** à taper.
- ✔️ **Le seul prérequis** : Un ordinateur sous **Windows 10 ou Windows 11**. L'installateur ajoute WebView2 si nécessaire.

---

## 🚀 Installation

1. [Téléchargez `Gestion-de-Stock-Setup.exe`](./Gestion-de-Stock-Setup.exe) puis lancez-le.
2. Suivez les étapes de l'installateur. L'installation est limitée à votre compte Windows et ne nécessite pas les droits administrateur.
3. Lancez **Gestion de Stock** depuis le raccourci créé sur le Bureau ou dans le menu Démarrer.

La base de données, les logos et les PDF de l'installation sont enregistrés dans votre dossier personnel Windows (`%APPDATA%\Gestion de Stock`) et ne sont pas inclus dans l'installateur. Ils restent ainsi séparés des fichiers du programme.

## 🔐 Accès administrateur et utilisateur

Au premier démarrage, créez un mot de passe administrateur et un mot de passe utilisateur différents (8 caractères minimum). Les mots de passe sont enregistrés sous forme de hachages sécurisés.

- **Administrateur** : accès à la caisse, aux factures et bons, aux produits, au stock et aux paramètres.
- **Utilisateur** : accès uniquement à la caisse et aux factures et bons.

La session est fermée à la fermeture de l'application. L'administrateur peut modifier les deux mots de passe depuis **Paramètres**.

---

## 🛠️ Que faire si Windows affiche un message au premier lancement ?

### 1. Écran bleu « Windows a protégé votre ordinateur » (SmartScreen)
C'est une alerte de sécurité standard de Windows lorsque vous recevez un nouveau fichier `.exe`.

Pour l'ouvrir en 2 clics :
1. Cliquez sur le lien **« Informations complémentaires »** sous le message.
2. Cliquez sur le bouton **« Exécuter quand même »** qui apparaît.  
*(L'application se lance et ne vous le redemandera plus).*

*Astuce alternative :* Clic droit sur `Gestion-de-Stock-Setup.exe` ➔ **Propriétés** ➔ tout en bas, cochez la case **« Débloquer »** ➔ **Appliquer** ➔ **OK**.

---

## 📦 Ce que fait le logiciel

- **🛒 Caisse, Douchette & TPE Bancaire** : Branchez votre lecteur code-barres USB, scannez vos articles, choisissez le mode de règlement. En sélectionnant **Carte Bancaire**, le logiciel dialogue directement avec votre terminal de paiement (TPE branché en USB Série ou Réseau Ethernet, ou via le simulateur intégré pour tests).
- **📦 Produits & Stocks** : Catalogue avec recherche instantanée, entrées de stock fournisseurs, sorties, casse, ajustements physiques d'inventaire et alertes automatiques de stock bas.
- **🧾 Factures & Bons de livraison** : Numérotation séquentielle officielle sans trou (`FAC-2026-XXXXX`), logo personnalisable, calcul automatique HT/TVA/TTC et génération de PDF instantanée.
- **⚙️ Paramètres & Sauvegarde** : Sauvegarde sécurisée de toutes vos données en un clic et configuration universelle du TPE.

---

## 💳 Terminal de Paiement (TPE) & Paiement Carte Bancaire

Le logiciel intègre un module universel de paiement par carte bancaire (**`Payement_carte`**) compatible avec tous les matériels :

- **🧪 Mode Simulateur (Par défaut pour vos tests) :**
  - Permet de tester immédiatement l'encaissement CB sans aucun TPE branché.
  - Dans **Paramètres**, vous pouvez choisir le scénario de test souhaité : *Paiement Validé*, *Paiement Refusé par la banque*, *Annulation sur le terminal*, ou *Aléatoire*.
- **🔌 Mode Câble USB / Série (Port COM) :**
  - Compatible avec le protocole standard français **Concert V3 / CB PR V3.0** (Ingenico Desk/Move, PAX, Verifone).
  - Détection automatique des ports COM dans **Paramètres** avec bouton de test de communication.
- **🌐 Mode Réseau Local (Ethernet / IP) :**
  - Pour les TPE reliés par câble réseau Ethernet ou Wi-Fi à votre box internet.

**Flux en caisse :**
1. Cliquez sur **Carte Bancaire** puis **Valider l'encaissement**.
2. Une modale interactive affiche le montant exact transmis au TPE et l'état en direct.
3. 🟢 **Validé** : La facture est validée et ouverte en PDF automatiquement.
4. 🔴 **Refusé** : Motif affiché clairement avec bouton **« Recommencer sur le TPE »** ou **« Payer en Espèces »**.
5. ⚪ **Annulé** : Possibilité de **« Réessayer le paiement CB »** ou de changer de moyen de paiement.

---

## 🌐 Utilisation sur 2 Ordinateurs (Multi-Postes avec Base Partagée)

Le logiciel gère nativement le fonctionnement sur 2 ordinateurs connectés au même réseau local (box Wi-Fi ou câble Ethernet) :

### 🖥️ Poste 1 (Principal / Administrateur) :
1. Lancez **`Gestion-de-Stock.exe`**.
2. Connectez-vous avec le compte **Administrateur**.
3. **Le serveur de partage démarre automatiquement en arrière-plan.**  
   En bas du menu de gauche, un badge vert affiche : `🟢 PC 1 : <Votre_IP>:8765` (ex: `192.168.1.45:8765`).
4. Si le pare-feu Windows demande une autorisation au premier lancement, cochez **Réseau privé** et cliquez sur **Autoriser l'accès**.

### 💻 Poste 2 (Secondaire / Caisse Utilisateur) :
1. Lancez **`Gestion-de-Stock.exe`** sur le deuxième PC.
2. En bas du menu de gauche, cliquez sur le bouton de connexion réseau 📡 (**« Lier ce poste au PC 1 »**).
3. Entrez l'adresse IP affichée sur le PC 1 (ex: `192.168.1.45`) puis cliquez sur **Tester & Lier**.
4. Le voyant passe au vert : le Poste 2 est désormais connecté et partage en temps réel la même base de données que le Poste 1 (produits, stocks, caisse, factures) !

---

## استخدام البرنامج على حاسوبين في نفس الشبكة

يمكن جعل حاسوب واحد هو الرئيسي؛ قاعدة البيانات تبقى عليه، ويتصل به الحاسوب الثاني عبر الشبكة المحلية. لا تضع ملف SQLite داخل مجلد شبكة مشترك ولا تفتح المنفذ على الإنترنت.

1. ثبّت النسخة المبنية حديثاً على الحاسوبين. اترك التطبيق مغلقاً أثناء إنشاء ملف الإعدادات.
2. على الحاسوب الرئيسي، افتح `%APPDATA%\Gestion de Stock` وأنشئ الملف `network.json` بهذا الشكل:
   ```json
   {
     "mode": "server",
     "listen": "0.0.0.0:8765",
     "token": "ضع_هنا_رمزاً_عشوائياً_طويلاً"
   }
   ```
   يجب ألا يقل الرمز عن 32 حرفاً. يمكنك توليد رمز في PowerShell بالأمر:
   ```powershell
   [guid]::NewGuid().ToString("N") + [guid]::NewGuid().ToString("N")
   ```
3. شغّل البرنامج على الحاسوب الرئيسي واتركه مفتوحاً. عند أول طلب من Windows Firewall، اسمح بالاتصال على **الشبكات الخاصة فقط**.
4. اعرف عنوان IPv4 للحاسوب الرئيسي باستخدام `ipconfig` (مثلاً `192.168.1.20`). على الحاسوب الثاني، أنشئ الملف نفسه في `%APPDATA%\Gestion de Stock`:
   ```json
   {
     "mode": "client",
     "serverUrl": "http://192.168.1.20:8765",
     "token": "استخدم_نفس_الرمز_الموجود_على_الرئيسي"
   }
   ```
5. شغّل التطبيق على الحاسوب الثاني. يجب أن يكون الحاسوبان على الشبكة نفسها، وأن يبقى الرئيسي مشغلاً. تغييرات المنتجات والمخزون والفواتير تُحفظ مباشرة في قاعدة بيانات الحاسوب الرئيسي.

تُنزّل ملفات PDF والنسخ الاحتياطية التي تُنشأ من الحاسوب الثاني على ذلك الحاسوب. يمكن اختيار شعار الشركة من الحاسوب الثاني وسيُحفظ على الرئيسي. الاتصال يستخدم HTTP ورمز مشاركة، لذا استخدمه على شبكة محلية موثوقة فقط، ولا تفتح المنفذ 8765 على موجّه الإنترنت. إذا تغيّر عنوان الحاسوب الرئيسي، حدّث `serverUrl` على الحاسوب الثاني.

---

## 💻 Annexe : Pour les développeurs uniquement (Optionnel)

<details>
<summary>Cliquez ici pour afficher les commandes de compilation depuis les sources</summary>

Si vous souhaitez recompiler l'application vous-même depuis le code source :

**Prérequis développeur :**
- Go >= 1.22
- Node.js >= 18 & npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- NSIS (pour générer l'installateur)

**Commandes :**
```bash
# Lancement en mode développement
wails dev

# Génération de l'installateur Windows, sans droits administrateur à l'installation
wails build -nsis -installscope user
```
</details>

---

## 📄 Licence

Copyright © 2026 **Adel Mahmoud et Sidou**. Tous droits réservés.
