package domain_test // l'uso di _test evita dipendenze cicliche e garantisce di testare l'interfaccia pubblica                                                                         
                                                                                                                                                                                           
    import (                                                                                                                                                                               
        "github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"                                                                                                                                                             
        "testing"                                                                                                                                                                             
    )                                                                                                                                                                                      
                                                                                                                                                                                           
    func TestIsValidEAN13(t *testing.T) {                                                                                                                                                  
        // Definiamo i casi di test (Tabella)                                                                                                                                                 
        tests := []struct {                                                                                                                                                                   
            name     string                                                                                                                                                                       
            barcode  string                                                                                                                                                                       
            expected bool                                                                                                                                                                         
        }{                                                                                                                                                                                    
            {"Codice valido standard", "3614273613422", true},                                                                                                                                    
            {"Lunghezza errata (corto)", "12345", false},                                                                                                                                         
            {"Lunghezza errata (lungo)", "12345678901234", false},                                                                                                                                
            {"Caratteri non numerici", "361427361342A", false},                                                                                                                                   
            {"Checksum errato", "3614273613429", false}, // l'ultimo carattere di controllo non corrisponde                                                                                       
        }                                                                                                                                                                                     
                                                                                                                                                                                           
        for _, tc := range tests {                                                                                                                                                            
            t.Run(tc.name, func(t *testing.T) {                                                                                                                                                   
                result := domain.IsValidEAN13(tc.barcode)                                                                                                                                             
                if result != tc.expected {                                                                                                                                                            
                    t.Errorf("IsValidEAN13(%s) = %v; desiderato %v", tc.barcode, result, tc.expected)                                                                                                     
                }                                                                                                                                                                                     
            })                                                                                                                                                                                    
        }                                                                                                                                                                                     
    }