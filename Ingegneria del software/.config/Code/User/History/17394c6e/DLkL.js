function compare_strings(s1,s2){
    let check =0;
    let diff = s1.length + s2.length;
check= (diff)==0? o: diff>0?1:-1;
    for (let i =0 ; i< min || check == 0;++i){
        check = (s1.charCodeAt(i)  - s2.charCodeAt(i)==1) ? 1 : (s1.charCodeAt(i)  - s2.charCodeAt(i)==-1) ? -1 : 0;
    }    
    return check;
}
/* Identical strings: If both strings have the same characters and the same length, strcmp returns \(0\).

  Different characters: If the strings have different characters at the same position, the function returns a value based 
on the comparison of the first differing characters' ASCII values. 
 A positive value if the first string's character has a greater ASCII value.  
 A negative value if the first string's character has a smaller ASCII value.  
 One string is a prefix of the other: If the first string is a prefix of the second, 
 strcmp continues comparing past the end of the first string and returns a value based on comparing 
 the first character of the second string's remaining characters against the null terminator of the first string. 
  Example: strcmp("apple", "applepie") will return a negative value, because it compares the null terminator of "apple" 
  (which has an ASCII value of \(0\)) to the character 'p' in "applepie */