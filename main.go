package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// NAME,STREET_LINE_1,STREET_LINE_2,CITY,STATE,ZIP,PROFESSION_JOB_TITLE,EMPLOYERS_NAME_SPECIFIC_FIELD,TRANSACTION_TYPE,COMMITTEE_NAME,COMMITTEE_SBOE_ID,COMMITTEE_STREET_1,COMMITTEE_STREET_2,COMMITTEE_CITY,COMMITTEE_STATE,COMMITTEE_ZIP_CODE,REPORT_NAME,DATE_OCCURED,ACCOUNT_CODE,AMOUNT,FORM_OF_PAYMENT,PURPOSE,CANDIDATE_REFERENDUM_NAME,DECLARATION,ORIGINAL_NAME,NAME_ID

type input struct {
	name                          string
	street_line_1                 string
	street_line_2                 string
	city                          string
	state                         string
	zip                           string
	profession_job_title          string
	employers_name_specific_field string
	transaction_type              string
	committee_name                string
	committee_sboe_id             string
	committee_street_1            string
	committee_street_2            string
	committee_city                string
	committee_state               string
	committee_zip_code            string
	report_name                   string
	date_occured                  string
	account_code                  string
	amount                        string
	form_of_payment               string
	purpose                       string
	candidate_referendum_name     string
	declaration                   string
	original_name                 string
	name_id                       string
}

type CSVAble interface {
	ToCSV() []string
	Header() []string
}

type transaction struct {
	source_contribution_id int
	contributor_id,
	transaction_type,
	committee_name,
	canon_committee_sboe_id,
	transaction_category,
	date_occurred,
	amount,
	form_of_payment,
	purpose,
	candidate_referendum_name,
	declaration,
	original_committee_sboe_id,
	original_account_id string
}

func (t transaction) ToCSV() []string {
	return []string{
		fmt.Sprint(t.source_contribution_id),
		t.contributor_id,
		t.transaction_type,
		t.committee_name,
		t.canon_committee_sboe_id,
		t.transaction_category,
		t.date_occurred,
		t.amount,
		t.form_of_payment,
		t.purpose,
		t.candidate_referendum_name,
		t.declaration,
		t.original_committee_sboe_id,
		t.original_account_id,
	}
}

func (t transaction) Header() []string {
	return []string{
		"source_contribution_id",
		"contributor_id",
		"transaction_type",
		"committee_name",
		"canon_committee_sboe_id",
		"transaction_category",
		"date_occurred",
		"amount",
		"form_of_payment",
		"purpose",
		"candidate_referendum_name",
		"declaration",
		"original_committee_sboe_id",
		"original_account_id",
	}
}

func (t *transaction) setId(id int) {
	t.source_contribution_id = id
}

type account struct {
	account_id    string // TODO: probably change this
	name          string
	street_line_1 string
	street_line_2 string
	city          string
	state         string
	zip           string
	profession    string
	employer_name string
}

func (a account) ToCSV() []string {
	return []string{
		a.account_id,
		a.name,
		a.street_line_1,
		a.street_line_2,
		a.city,
		a.state,
		a.zip,
		a.profession,
		a.employer_name,
	}
}

func (a account) Header() []string {
	return []string{
		"account_id",
		"name",
		"street_line_1",
		"street_line_2",
		"city",
		"state",
		"zip",
		"profession",
		"employer_name",
	}
}

type committee struct {
	sboe_id            string
	committee_name     string
	committee_street_1 string
	committee_street_2 string
	committee_city     string
	committee_state    string
	committee_zip_code string
	candidate_full_name,
	candidate_first_last_name,
	party,
	office,
	juris string
}

func (c committee) ToCSV() []string {
	return []string{
		c.sboe_id,
		c.committee_name,
		c.committee_street_1,
		c.committee_street_2,
		c.committee_city,
		c.committee_state,
		c.committee_zip_code,
		c.candidate_full_name,
		c.candidate_first_last_name,
		c.party,
		c.office,
		c.juris,
	}
}

func (c committee) Header() []string {
	return []string{
		"sboe_id",
		"committee_name",
		"committee_street_1",
		"committee_street_2",
		"committee_city",
		"committee_state",
		"committee_zip_code",
		"candidate_full_name",
		"candidate_first_last_name",
		"party",
		"office",
		"juris",
	}
}

func (c *committee) addReferenceInfo(m map[string]referenceCommittee) {
	ref, ok := m[c.sboe_id]
	if !ok {
		fmt.Println("no reference found for", c.sboe_id)
		return
	}
	c.candidate_first_last_name = ref.candidate_first_last_name
	c.candidate_full_name = ref.candidate_full_name
	c.party = ref.party
	c.office = ref.office
	c.juris = ref.juris

}

type referenceCommittee struct {
	sboe_id,
	committee_name,
	committee_type,
	committee_street_1,
	committee_street_2,
	committee_city,
	committee_state,
	committee_full_zip,
	candidate_first_name,
	candidate_middle_name,
	candidate_last_name,
	candidate_full_name,
	candidate_first_last_name,
	treasurer_first_name,
	treasurer_middle_name,
	treasurer_last_name,
	treasurer_email,
	asst_treasurer_first_name,
	asst_treasurer_middle_name,
	asst_treasurer_last_name,
	asst_treasurer_email,
	treasurer_street_1,
	treasurer_street_2,
	treasurer_city,
	treasurer_state,
	treasurer_full_zip,
	party,
	office,
	juris string
}

func buildReferenceCommitteeMap(path string) (map[string]referenceCommittee, error) {
	m := map[string]referenceCommittee{}
	file, err := os.Open(path)
	if err != nil {
		return m, err
	}
	r := csv.NewReader(file)
	// skip header
	if _, err := r.Read(); err != nil {
		return m, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record = cleanLine(record)
		c := referenceCommittee{
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
			record[7],
			record[8],
			record[9],
			record[10],
			record[11],
			record[12],
			record[13],
			record[14],
			record[15],
			record[16],
			record[17],
			record[18],
			record[19],
			record[20],
			record[21],
			record[22],
			record[23],
			record[24],
			record[25],
			record[26],
			record[27],
			record[28],
		}
		m[c.sboe_id] = c

	}

	return m, nil
}

func newInput(in []string) input {
	return input{
		in[0],
		in[1],
		in[2],
		in[3],
		in[4],
		in[5],
		in[6],
		in[7],
		in[8],
		in[9],
		in[10],
		in[11],
		in[12],
		in[13],
		in[14],
		in[15],
		in[16],
		in[17],
		in[18],
		in[19],
		in[20],
		in[21],
		in[22],
		in[23],
		in[24],
		in[25],
	}
}

func cleanLine(in []string) []string {
	out := make([]string, len(in))

	for i, x := range in {
		out[i] = strings.Trim(x, " ")
		if out[i] == "NULL" { // This could cause issues for people named "Null"
			out[i] = ""
		}
	}
	return out
}

func inputToTypes(in input) (transaction, account, committee) {
	t := transaction{
		contributor_id:             in.name_id,
		original_committee_sboe_id: in.committee_sboe_id,
		original_account_id:        in.name_id,
		// TODO: delete one of these
		transaction_type:          in.transaction_type,
		transaction_category:      in.transaction_type,
		committee_name:            in.committee_name,
		canon_committee_sboe_id:   in.committee_sboe_id,
		date_occurred:             in.date_occured,
		amount:                    in.amount,
		form_of_payment:           in.form_of_payment,
		purpose:                   in.purpose,
		candidate_referendum_name: in.candidate_referendum_name,
		declaration:               in.declaration,
	}

	a := account{
		account_id:    in.name_id,
		name:          in.name,
		street_line_1: in.street_line_1,
		street_line_2: in.street_line_2,
		city:          in.city,
		state:         in.state,
		zip:           in.zip,
		profession:    in.profession_job_title,
		employer_name: in.employers_name_specific_field,
	}

	c := committee{
		sboe_id:            in.committee_sboe_id,
		committee_name:     in.committee_name,
		committee_street_1: in.committee_street_1,
		committee_street_2: in.committee_street_2,
		committee_city:     in.city,
		committee_state:    in.state,
		committee_zip_code: in.zip,
		// candidate_full_name,
		// candidate_first_last_name,
		// party,
		// office,
		// juris
	}

	return t, a, c
}

func buildCSVWriter(fileName string, c CSVAble) (*csv.Writer, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	err = writer.Write(c.Header())
	if err != nil {
		return nil, err
	}
	return writer, nil
}

func main() {
	in, err := os.Open("./input.csv")
	if err != nil {
		log.Fatal(err)
	}
	transactionWriter, err := buildCSVWriter("./transactions.csv", transaction{})
	defer transactionWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
	accountWriter, err := buildCSVWriter("./accounts.csv", account{})
	defer accountWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
	committeeWriter, err := buildCSVWriter("./committees.csv", committee{})
	defer committeeWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}

	m, err := buildReferenceCommitteeMap("./reference_committees.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(in)
	// skip header
	if _, err := r.Read(); err != nil {
		log.Fatal(err)
	}

	// limit := 5
	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("processing record", count+1)

		t, a, c := inputToTypes(newInput(cleanLine(record)))
		t.setId(count)
		err = transactionWriter.Write(t.ToCSV())
		if err != nil {
			log.Fatal(err)
		}
		err = accountWriter.Write(a.ToCSV())
		if err != nil {
			log.Fatal(err)
		}

		c.addReferenceInfo(m)
		err = committeeWriter.Write(c.ToCSV())
		if err != nil {
			log.Fatal(err)
		}
		// if count > limit {
		// 	break
		// }
		count++
	}

}

// TODO: add id and report_name to transactions
// ensure out csvs match expected input/update schema if necessary
// may need to quote csv output files
// https://stackoverflow.com/questions/21324133/always-quote-csv-values
// test import new files
// ???
// profit
