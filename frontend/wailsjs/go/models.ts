export namespace factures {
	
	export class Bon {
	    id: string;
	    typeBon: string;
	    numero: string;
	    sequenceNumero: number;
	    statut: string;
	    dateBon: string;
	    referenceExterne: string;
	    factureId?: string;
	    entrepriseDenomination: string;
	    entrepriseAdresse: string;
	    entrepriseTelephone: string;
	    entrepriseEmail: string;
	    clientId?: string;
	    clientNom: string;
	    clientAdresse: string;
	    clientCodePostal: string;
	    clientVille: string;
	    clientTelephone: string;
	    totalHt: number;
	    totalTva: number;
	    totalTtc: number;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Bon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.typeBon = source["typeBon"];
	        this.numero = source["numero"];
	        this.sequenceNumero = source["sequenceNumero"];
	        this.statut = source["statut"];
	        this.dateBon = source["dateBon"];
	        this.referenceExterne = source["referenceExterne"];
	        this.factureId = source["factureId"];
	        this.entrepriseDenomination = source["entrepriseDenomination"];
	        this.entrepriseAdresse = source["entrepriseAdresse"];
	        this.entrepriseTelephone = source["entrepriseTelephone"];
	        this.entrepriseEmail = source["entrepriseEmail"];
	        this.clientId = source["clientId"];
	        this.clientNom = source["clientNom"];
	        this.clientAdresse = source["clientAdresse"];
	        this.clientCodePostal = source["clientCodePostal"];
	        this.clientVille = source["clientVille"];
	        this.clientTelephone = source["clientTelephone"];
	        this.totalHt = source["totalHt"];
	        this.totalTva = source["totalTva"];
	        this.totalTtc = source["totalTtc"];
	        this.notes = source["notes"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LigneBon {
	    id: string;
	    bonId: string;
	    produitId?: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    quantite: number;
	    prixUnitaireHt: number;
	    tauxTva: number;
	    montantHt: number;
	    montantTva: number;
	    montantTtc: number;
	    ordre: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new LigneBon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.bonId = source["bonId"];
	        this.produitId = source["produitId"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.quantite = source["quantite"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.tauxTva = source["tauxTva"];
	        this.montantHt = source["montantHt"];
	        this.montantTva = source["montantTva"];
	        this.montantTtc = source["montantTtc"];
	        this.ordre = source["ordre"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DetailBon {
	    bon: Bon;
	    lignes: LigneBon[];
	
	    static createFrom(source: any = {}) {
	        return new DetailBon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bon = this.convertValues(source["bon"], Bon);
	        this.lignes = this.convertValues(source["lignes"], LigneBon);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LigneFacture {
	    id: string;
	    factureId: string;
	    produitId?: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    quantite: number;
	    prixUnitaireHt: number;
	    tauxTva: number;
	    montantHt: number;
	    montantTva: number;
	    montantTtc: number;
	    ordre: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new LigneFacture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.factureId = source["factureId"];
	        this.produitId = source["produitId"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.quantite = source["quantite"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.tauxTva = source["tauxTva"];
	        this.montantHt = source["montantHt"];
	        this.montantTva = source["montantTva"];
	        this.montantTtc = source["montantTtc"];
	        this.ordre = source["ordre"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Facture {
	    id: string;
	    numero: string;
	    sequenceNumero: number;
	    statut: string;
	    modeReglement: string;
	    dateFacture: string;
	    dateEcheance: string;
	    conditionsReglement: string;
	    entrepriseDenomination: string;
	    entrepriseFormeJuridique: string;
	    entrepriseCapital: number;
	    entrepriseAdresse: string;
	    entrepriseCodePostal: string;
	    entrepriseVille: string;
	    entrepriseSiren: string;
	    entrepriseSiret: string;
	    entrepriseRcs: string;
	    entrepriseNumeroTva: string;
	    entrepriseTelephone: string;
	    entrepriseEmail: string;
	    mentionsLegales: string;
	    devise: string;
	    clientId?: string;
	    clientNom: string;
	    clientAdresse: string;
	    clientCodePostal: string;
	    clientVille: string;
	    clientTelephone: string;
	    clientEmail: string;
	    clientSirenSiret: string;
	    clientNumeroTva: string;
	    totalHt: number;
	    totalTva: number;
	    totalTtc: number;
	    montantRegle: number;
	    resteAPayer: number;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Facture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.numero = source["numero"];
	        this.sequenceNumero = source["sequenceNumero"];
	        this.statut = source["statut"];
	        this.modeReglement = source["modeReglement"];
	        this.dateFacture = source["dateFacture"];
	        this.dateEcheance = source["dateEcheance"];
	        this.conditionsReglement = source["conditionsReglement"];
	        this.entrepriseDenomination = source["entrepriseDenomination"];
	        this.entrepriseFormeJuridique = source["entrepriseFormeJuridique"];
	        this.entrepriseCapital = source["entrepriseCapital"];
	        this.entrepriseAdresse = source["entrepriseAdresse"];
	        this.entrepriseCodePostal = source["entrepriseCodePostal"];
	        this.entrepriseVille = source["entrepriseVille"];
	        this.entrepriseSiren = source["entrepriseSiren"];
	        this.entrepriseSiret = source["entrepriseSiret"];
	        this.entrepriseRcs = source["entrepriseRcs"];
	        this.entrepriseNumeroTva = source["entrepriseNumeroTva"];
	        this.entrepriseTelephone = source["entrepriseTelephone"];
	        this.entrepriseEmail = source["entrepriseEmail"];
	        this.mentionsLegales = source["mentionsLegales"];
	        this.devise = source["devise"];
	        this.clientId = source["clientId"];
	        this.clientNom = source["clientNom"];
	        this.clientAdresse = source["clientAdresse"];
	        this.clientCodePostal = source["clientCodePostal"];
	        this.clientVille = source["clientVille"];
	        this.clientTelephone = source["clientTelephone"];
	        this.clientEmail = source["clientEmail"];
	        this.clientSirenSiret = source["clientSirenSiret"];
	        this.clientNumeroTva = source["clientNumeroTva"];
	        this.totalHt = source["totalHt"];
	        this.totalTva = source["totalTva"];
	        this.totalTtc = source["totalTtc"];
	        this.montantRegle = source["montantRegle"];
	        this.resteAPayer = source["resteAPayer"];
	        this.notes = source["notes"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DetailFacture {
	    facture: Facture;
	    lignes: LigneFacture[];
	
	    static createFrom(source: any = {}) {
	        return new DetailFacture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.facture = this.convertValues(source["facture"], Facture);
	        this.lignes = this.convertValues(source["lignes"], LigneFacture);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FiltreBons {
	    typeBon?: string;
	    recherche: string;
	    statut: string;
	    clientId: string;
	    dateDebut: string;
	    dateFin: string;
	    page: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new FiltreBons(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.typeBon = source["typeBon"];
	        this.recherche = source["recherche"];
	        this.statut = source["statut"];
	        this.clientId = source["clientId"];
	        this.dateDebut = source["dateDebut"];
	        this.dateFin = source["dateFin"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}
	export class FiltreFactures {
	    recherche: string;
	    statut?: string;
	    clientId: string;
	    dateDebut: string;
	    dateFin: string;
	    page: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new FiltreFactures(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recherche = source["recherche"];
	        this.statut = source["statut"];
	        this.clientId = source["clientId"];
	        this.dateDebut = source["dateDebut"];
	        this.dateFin = source["dateFin"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}
	
	
	export class SaisieLigneBon {
	    produitId?: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    quantite: number;
	    prixUnitaireHt: number;
	    tauxTva: number;
	
	    static createFrom(source: any = {}) {
	        return new SaisieLigneBon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.quantite = source["quantite"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.tauxTva = source["tauxTva"];
	    }
	}
	export class NouveauBon {
	    typeBon: string;
	    dateBon: string;
	    referenceExterne: string;
	    clientId?: string;
	    clientNom: string;
	    notes: string;
	    lignes: SaisieLigneBon[];
	
	    static createFrom(source: any = {}) {
	        return new NouveauBon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.typeBon = source["typeBon"];
	        this.dateBon = source["dateBon"];
	        this.referenceExterne = source["referenceExterne"];
	        this.clientId = source["clientId"];
	        this.clientNom = source["clientNom"];
	        this.notes = source["notes"];
	        this.lignes = this.convertValues(source["lignes"], SaisieLigneBon);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SaisieLigneFacture {
	    produitId?: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    quantite: number;
	    prixUnitaireHt: number;
	    tauxTva: number;
	
	    static createFrom(source: any = {}) {
	        return new SaisieLigneFacture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.quantite = source["quantite"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.tauxTva = source["tauxTva"];
	    }
	}
	export class NouvelleFacture {
	    clientId?: string;
	    clientNom: string;
	    clientAdresse: string;
	    clientSirenSiret: string;
	    modeReglement: string;
	    dateFacture: string;
	    dateEcheance: string;
	    conditionsReglement: string;
	    montantRegle: number;
	    notes: string;
	    lignes: SaisieLigneFacture[];
	
	    static createFrom(source: any = {}) {
	        return new NouvelleFacture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.clientNom = source["clientNom"];
	        this.clientAdresse = source["clientAdresse"];
	        this.clientSirenSiret = source["clientSirenSiret"];
	        this.modeReglement = source["modeReglement"];
	        this.dateFacture = source["dateFacture"];
	        this.dateEcheance = source["dateEcheance"];
	        this.conditionsReglement = source["conditionsReglement"];
	        this.montantRegle = source["montantRegle"];
	        this.notes = source["notes"];
	        this.lignes = this.convertValues(source["lignes"], SaisieLigneFacture);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageBons {
	    bons: Bon[];
	    total: number;
	    page: number;
	    limit: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PageBons(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bons = this.convertValues(source["bons"], Bon);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.totalPages = source["totalPages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageFactures {
	    factures: Facture[];
	    total: number;
	    page: number;
	    limit: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PageFactures(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.factures = this.convertValues(source["factures"], Facture);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.totalPages = source["totalPages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace parametres {
	
	export class ModificationParametres {
	    denomination: string;
	    formeJuridique: string;
	    capital: number;
	    adresse: string;
	    codePostal: string;
	    ville: string;
	    siren: string;
	    siret: string;
	    rcs: string;
	    numeroTva: string;
	    telephone: string;
	    email: string;
	    siteWeb: string;
	    logoPath: string;
	    devise: string;
	    tauxTvaDefaut: number;
	    tauxTvaDisponibles: string;
	    mentionsLegales: string;
	    conditionsReglement: string;
	    formatNumerotationFacture: string;
	    formatNumerotationBon: string;
	
	    static createFrom(source: any = {}) {
	        return new ModificationParametres(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.denomination = source["denomination"];
	        this.formeJuridique = source["formeJuridique"];
	        this.capital = source["capital"];
	        this.adresse = source["adresse"];
	        this.codePostal = source["codePostal"];
	        this.ville = source["ville"];
	        this.siren = source["siren"];
	        this.siret = source["siret"];
	        this.rcs = source["rcs"];
	        this.numeroTva = source["numeroTva"];
	        this.telephone = source["telephone"];
	        this.email = source["email"];
	        this.siteWeb = source["siteWeb"];
	        this.logoPath = source["logoPath"];
	        this.devise = source["devise"];
	        this.tauxTvaDefaut = source["tauxTvaDefaut"];
	        this.tauxTvaDisponibles = source["tauxTvaDisponibles"];
	        this.mentionsLegales = source["mentionsLegales"];
	        this.conditionsReglement = source["conditionsReglement"];
	        this.formatNumerotationFacture = source["formatNumerotationFacture"];
	        this.formatNumerotationBon = source["formatNumerotationBon"];
	    }
	}
	export class ParametresEntreprise {
	    id: string;
	    denomination: string;
	    formeJuridique: string;
	    capital: number;
	    adresse: string;
	    codePostal: string;
	    ville: string;
	    siren: string;
	    siret: string;
	    rcs: string;
	    numeroTva: string;
	    telephone: string;
	    email: string;
	    siteWeb: string;
	    logoPath: string;
	    devise: string;
	    tauxTvaDefaut: number;
	    tauxTvaDisponibles: string;
	    mentionsLegales: string;
	    conditionsReglement: string;
	    formatNumerotationFacture: string;
	    formatNumerotationBon: string;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ParametresEntreprise(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.denomination = source["denomination"];
	        this.formeJuridique = source["formeJuridique"];
	        this.capital = source["capital"];
	        this.adresse = source["adresse"];
	        this.codePostal = source["codePostal"];
	        this.ville = source["ville"];
	        this.siren = source["siren"];
	        this.siret = source["siret"];
	        this.rcs = source["rcs"];
	        this.numeroTva = source["numeroTva"];
	        this.telephone = source["telephone"];
	        this.email = source["email"];
	        this.siteWeb = source["siteWeb"];
	        this.logoPath = source["logoPath"];
	        this.devise = source["devise"];
	        this.tauxTvaDefaut = source["tauxTvaDefaut"];
	        this.tauxTvaDisponibles = source["tauxTvaDisponibles"];
	        this.mentionsLegales = source["mentionsLegales"];
	        this.conditionsReglement = source["conditionsReglement"];
	        this.formatNumerotationFacture = source["formatNumerotationFacture"];
	        this.formatNumerotationBon = source["formatNumerotationBon"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace produits {
	
	export class FiltreProduits {
	    recherche: string;
	    alerteStock: boolean;
	    actif?: boolean;
	    page: number;
	    limit: number;
	    triChamp: string;
	    triAscendant: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FiltreProduits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recherche = source["recherche"];
	        this.alerteStock = source["alerteStock"];
	        this.actif = source["actif"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.triChamp = source["triChamp"];
	        this.triAscendant = source["triAscendant"];
	    }
	}
	export class ModificationProduit {
	    id: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockMinimum: number;
	    unite: string;
	    actif: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModificationProduit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	        this.actif = source["actif"];
	    }
	}
	export class NouveauProduit {
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockInitial: number;
	    stockMinimum: number;
	    unite: string;
	
	    static createFrom(source: any = {}) {
	        return new NouveauProduit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockInitial = source["stockInitial"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	    }
	}
	export class Produit {
	    id: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockActuel: number;
	    stockMinimum: number;
	    unite: string;
	    actif: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Produit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockActuel = source["stockActuel"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	        this.actif = source["actif"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageProduits {
	    produits: Produit[];
	    total: number;
	    page: number;
	    limit: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PageProduits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produits = this.convertValues(source["produits"], Produit);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.totalPages = source["totalPages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace stock {
	
	export class AlerteStock {
	    produitId: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    stockActuel: number;
	    stockMinimum: number;
	    unite: string;
	
	    static createFrom(source: any = {}) {
	        return new AlerteStock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.stockActuel = source["stockActuel"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	    }
	}
	export class FiltreMouvements {
	    produitId: string;
	    typeMouvement?: string;
	    dateDebut: string;
	    dateFin: string;
	    page: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new FiltreMouvements(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.typeMouvement = source["typeMouvement"];
	        this.dateDebut = source["dateDebut"];
	        this.dateFin = source["dateFin"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}
	export class MouvementDetail {
	    id: string;
	    produitId: string;
	    typeMouvement: string;
	    quantite: number;
	    stockAvant: number;
	    stockApres: number;
	    prixUnitaireHt: number;
	    referenceDocument: string;
	    motif: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	
	    static createFrom(source: any = {}) {
	        return new MouvementDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.produitId = source["produitId"];
	        this.typeMouvement = source["typeMouvement"];
	        this.quantite = source["quantite"];
	        this.stockAvant = source["stockAvant"];
	        this.stockApres = source["stockApres"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.referenceDocument = source["referenceDocument"];
	        this.motif = source["motif"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MouvementStock {
	    id: string;
	    produitId: string;
	    typeMouvement: string;
	    quantite: number;
	    stockAvant: number;
	    stockApres: number;
	    prixUnitaireHt: number;
	    referenceDocument: string;
	    motif: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new MouvementStock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.produitId = source["produitId"];
	        this.typeMouvement = source["typeMouvement"];
	        this.quantite = source["quantite"];
	        this.stockAvant = source["stockAvant"];
	        this.stockApres = source["stockApres"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.referenceDocument = source["referenceDocument"];
	        this.motif = source["motif"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageMouvements {
	    mouvements: MouvementDetail[];
	    total: number;
	    page: number;
	    limit: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PageMouvements(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mouvements = this.convertValues(source["mouvements"], MouvementDetail);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.totalPages = source["totalPages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SaisieAjustement {
	    produitId: string;
	    nouveauStock: number;
	    motif: string;
	
	    static createFrom(source: any = {}) {
	        return new SaisieAjustement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.nouveauStock = source["nouveauStock"];
	        this.motif = source["motif"];
	    }
	}
	export class SaisieMouvement {
	    produitId: string;
	    typeMouvement: string;
	    quantite: number;
	    prixUnitaireHt: number;
	    referenceDocument: string;
	    motif: string;
	
	    static createFrom(source: any = {}) {
	        return new SaisieMouvement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produitId = source["produitId"];
	        this.typeMouvement = source["typeMouvement"];
	        this.quantite = source["quantite"];
	        this.prixUnitaireHt = source["prixUnitaireHt"];
	        this.referenceDocument = source["referenceDocument"];
	        this.motif = source["motif"];
	    }
	}

}

