insert into public.mcp_blocks (id, name, description) values
  (1, 'Execution Drive', 'Initiative, accountability, and delivery focus'),
  (2, 'Cognitive Processing', 'Thinking, analysis, and decision-making style'),
  (3, 'Professional Values', 'Professional principles and motivations'),
  (4, 'Emotional Resilience', 'Response to pressure and setbacks'),
  (5, 'Learning Velocity', 'Learning and adaptability'),
  (6, 'Collaboration Depth', 'Teamwork and communication'),
  (7, 'Operational Rhythm', 'Pace, structure, and execution'),
  (8, 'Team Positioning', 'Role and contribution within a team')
on conflict (id) do update set name = excluded.name, description = excluded.description;

insert into public.mcp_questions (block_id, dimension_index, question_text, scale)
select
  block_id,
  dimension_index,
  format('I demonstrate %s consistently in my work.', lower(block_name)),
  '0-1'
from (
  select b.id as block_id, b.name as block_name, ((b.id - 1) * 8) + item as dimension_index
  from public.mcp_blocks b cross join generate_series(0, 7) as item
) dimensions
on conflict (dimension_index) do nothing;
