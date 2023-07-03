package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type proficiencyLevel = string

const (
    junior proficiencyLevel = "junior"
    expert proficiencyLevel = "expert"
    senior proficiencyLevel = "senior"
    teamLead proficiencyLevel = "team_lead"
)

const (
    employees_file_path = "./data/employees.csv"
    working_hours_file_path = "./data/working_hours.csv"
    teams_file_path = "./data/teams.csv"
    salary_configs_file_path = "./data/salary_configs.csv"
)

type employee struct {
    id int
    name string
    age int
    proficiencyLevel proficiencyLevel
    team *team
    workingHours [30][24]bool
}

type team struct {
    id int
    headMember *employee
    bonusMinWorkingHours int
    members []*employee
}

type salaryConfig struct {
    baseSalary int
    salaryPerHour int
    salaryPerExtraHour int
    officialWorkingHours int
    taxPercentage int
}

type database struct {
    employees []*employee
    teams []*team
    salaryConfigs map[proficiencyLevel]*salaryConfig
}

func main() {
    db := NewDatabase()
    
    ProcessInput(db)
}

func Shout(err error) {
    if err != nil {
        fmt.Println(err.Error())
    }
}

func NewDatabase() database {
    db := database{}
    db.salaryConfigs = make(map[string]*salaryConfig)

    err := db.readEmployeesFile()
    Shout(err)
    err = db.readWorkingHoursFile()
    Shout(err)
    err = db.readTeamsFile()
    Shout(err)
    err = db.readSalaryConfigsFile()
    Shout(err)

    return db
}

func ProcessInput(db database) {
    var input string
    for {
        fmt.Scan(&input)
        switch input {
        case "report_salaries":
            db.ReportSalaries()
        }
    }
}

func (db database) ReportSalaries() {
    fmt.Println("Meow")
}

func (db *database) readEmployeesFile() error {
    lines, err := readFile(employees_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        id, _ := strconv.Atoi(rows[0])
        age, _ := strconv.Atoi(rows[2])
        
        db.addEmployee(id, rows[1], age, rows[3])
    }

    return err
}

func (db *database) readWorkingHoursFile() error {
    lines, err := readFile(working_hours_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        employeeId, _ := strconv.Atoi(rows[0])
        day, _ := strconv.Atoi(rows[1])
        workingInterval := strings.Split(rows[2], "-")
        startTime, _ := strconv.Atoi(workingInterval[0])
        endTime, _ := strconv.Atoi(workingInterval[1])
        
        db.setWorkingHours(employeeId, day, startTime, endTime)
    }

    return err
}

func (db *database) readTeamsFile() error {
    lines, err := readFile(teams_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        teamId, _ := strconv.Atoi(rows[0])
        headMemberId, _ := strconv.Atoi(rows[1])
        membersIds := StringsToInts(strings.Split(rows[2], "$"))
        bonusMinWorkingHours, _ := strconv.Atoi(rows[3])
        
        db.addTeam(teamId, headMemberId, membersIds, bonusMinWorkingHours)
    }

    return err
}

func (db *database) readSalaryConfigsFile() error {
    lines, err := readFile(salary_configs_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        baseSalary, _ := strconv.Atoi(rows[1])
        salaryPerHour, _ := strconv.Atoi(rows[2])
        salaryPerExtraHour, _ := strconv.Atoi(rows[3])
        officialWorkingHours, _ := strconv.Atoi(rows[4])
        taxPercentage, _ := strconv.Atoi(rows[5])
        
        db.addSalaryConfig(rows[0], baseSalary, salaryPerHour, salaryPerExtraHour, officialWorkingHours, taxPercentage)
    }

    return err
}

func readFile(filepath string) ([]string, error) {
    var lines []string
    file, err := os.Open(filepath)
    if err != nil {
        return lines, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    if err := scanner.Err(); err != nil {
        return lines, err
    }

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, err
}

func (db *database) addEmployee(id int, name string, age int, level proficiencyLevel) {
    newEmployee := employee{id: id, name: name, age: age, proficiencyLevel: level}
    db.employees = append(db.employees, &newEmployee)
}

func (db *database) setWorkingHours(employeeId, day, startTime, endTime int) {
    employee, _ := db.getEmployee(employeeId)
    for hour := startTime; hour < endTime; hour++ {
        employee.workingHours[day-1][hour] = true
    }
}

func (db *database) addTeam(teamId, headMemberId int, membersIds []int, bonusMinWorkingHours int) {
    newTeam := team{id: teamId, bonusMinWorkingHours: bonusMinWorkingHours}

    headMember, _ := db.getEmployee(headMemberId)
    newTeam.headMember = headMember

    for _, memberId := range membersIds {
        emp, _ := db.getEmployee(memberId)
        emp.team = &newTeam
        newTeam.members = append(newTeam.members, emp)
    }

    db.teams = append(db.teams, &newTeam)
}

func (db *database) addSalaryConfig(level proficiencyLevel, baseSalary, salaryPerHour, salaryPerExtraHour, officialWorkingHours, taxPercentage int) {
    newSalaryConfig := salaryConfig{baseSalary: baseSalary, salaryPerHour: salaryPerHour, salaryPerExtraHour: salaryPerExtraHour, officialWorkingHours: officialWorkingHours, taxPercentage: taxPercentage}
    db.salaryConfigs[level] = &newSalaryConfig
}

func (db database) getEmployee(id int) (*employee, error) {
    for _, employee := range db.employees {
        if employee.id == id {
            return employee, nil
        }
    }

    return nil, errors.New("Employee Not Found")
}

func (emp employee) GetTotalWorkingHours() int {
    counter := 0
    fmt.Println(emp.workingHours[15])
    for _, day := range emp.workingHours {
        for _, hour := range day {
            if hour {
                counter++
            }
        }
    }

    return counter
}

func StringsToInts(strings []string) []int {
    res := []int{}
    for _, str := range strings {
        num, _ := strconv.Atoi(str)
        res = append(res, num)
    }
    return res
}
