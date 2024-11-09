from ratemyprofessor import School, get_professor_by_school_and_name, get_professors_by_school_and_name

def get_professor_url(name: str) -> str:
    ccny = School(224)
    professors = get_professors_by_school_and_name(ccny, name)
    if not professors:
        return None
    print(professors[0])
    return_url = f"https://www.ratemyprofessors.com/professor/{professors[0].id}"
    print(return_url)
    
    
if __name__ == "__main__":
    get_professor_url("grimmelmann")