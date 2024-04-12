# Explanation of Fuzzy-Matching Approach

Typos, Optical Character Recognition (OCR) errors, and spelling variations in
names are common challenges in biology. We implemented fuzzy-matching to
address these issues. However, fuzzy-matching requires a careful balance
between recall (minimizing false negatives) and precision (minimizing false
positives). We developed several strategies to achieve this balance.

## Fuzzy-Matching Rules

1. **Removal of Specific Epithet Suffixes**

Scientific names sometimes have multiple suffixes due to discrepancies
between the gender of the genus and specific epithet. These discrepancies are
eventually resolved, leaving names with variations. GNverifier matches names
using "stemmed canonical forms" where suffixes of specific epithets are
removed.

   **Example:**

   * **Name:** Adiantum davidii Franchet var. longispina Ching
   * **Stemmed canonical:** Adiantum david longispin
   * **Match:** Adiantum davidii var. longispinum

2. **Limiting Edit Distance to 1**

Edit distance is the minimum number of changes needed to transform one
string into another. We found that edit distances greater than 1 produce more
false positives than accurate results. Therefore, we set the default edit
distance to 1 when comparing "stemmed" strings.

   **Example:**

   * **Name:** Odiantum davidii Franchet var. longispina Ching
   * **Input "stemmed canonical":** Odiantum dauid longispin
   * **Database "stemmed canonical":** Adiantum dauid longispin
   * **Match:** Adiantum davidii var. longispinum 

The edit distance between stemmed canonicals is 1, so the result is
accepted. The "final" edit distance is calculated as 3 due to suffix
differences.

3. **Disregarding Differences in Short Words**

Shorter words have a higher probability of naturally matching other short words
with an edit distance of 1. We determined that words shorter than 5 characters
should not be fuzzy-matched to prevent inaccuracies. 

4. **No Fuzzy-Matching for Uninomial Names**

Uninomial scientific names might also naturally match other uninomial words
with an edit distance of 1. Therefore, we do not apply fuzzy-matching to
uninomial names.

## Options to Change Fuzzy-Matching Rules

The following options can significantly increase false positives and should only be used when results are additionally restricted or manually checked. Refer to the main documentation for instructions on setting these options for the Web User Interface and command-line program.

1. **Allowing Fuzzy-Matching of Uninomials**

Enables fuzzy-matching for uninomial names.


2. **Relaxing Fuzzy-Matching Rules**

Increases the maximum allowed edit distance to 2 and removes restrictions on
word length. Due to increased computation and potential for false positives,
this option limits the maximum number of input names to 50 (instead of 5000).
